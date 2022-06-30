package database

import (
	//consensusAPI "ether-rbbc/consensus/api"
	dbAPI "ether-rbbc/database/api"
	"github.com/ethereum/go-ethereum/cmd/utils"

	//"ether-rbbc/dbft"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/eth/downloader"
	"github.com/ethereum/go-ethereum/event"
	tmtLog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
	"time"

	"ether-rbbc/database/metrics"
	"ether-rbbc/database/web3"
)

const (
	txChanSize = 8192
)

// Database manages the underlying ethereum state for storage and processing
// and maintains the connection to Tendermint for forwarding txs.

// Database handles the chain database and VM.
type Database struct {
	eth    *eth.Ethereum
	ethCfg *eth.Config

	ethTxSub event.Subscription
	ethTxsCh chan core.NewTxsEvent

	ethState *EthState

	//consAPI consensusAPI.API
	logger  tmtLog.Logger
	metrics metrics.Metrics

	proposeChan chan *ethTypes.Transaction
	start       bool
}

var db *Database
//var ReceiveChan chan *ethTypes.Transaction
func NewDatabase(ctx *node.Node, ethCfg *eth.Config, logger tmtLog.Logger, metrics metrics.Metrics, proposeChan chan *ethTypes.Transaction) (*Database, error) {
	//ethereum, err := eth.New(ctx, ethCfg)
	_, ethereum := utils.RegisterEthService(ctx, ethCfg)
	/*if err != nil {
		return nil, err
	}*/

	currentBlock := ethereum.BlockChain().CurrentBlock()
	ethereum.EventMux().Post(core.ChainHeadEvent{currentBlock})

	ethereum.BlockChain().SetValidator(NullBlockValidator{})

	db = &Database{
		eth:         ethereum,
		ethCfg:      ethCfg,
		ethState:    NewEthState(ethereum, ethCfg, logger),
		logger:      logger,
		metrics:     metrics,
		proposeChan: proposeChan,
		start:       true,
	}

	return db, nil
}

func (db *Database) Ethereum() *eth.Ethereum {
	return db.eth
}

func (db *Database) Config() *eth.Config {
	return db.ethCfg
}

// ExecuteTx appends a transaction to the current block.
func (db *Database) ExecuteTx(tx *ethTypes.Transaction) error {
	db.logger.Info("Executing DB TX", "hash", tx.Hash().Hex(), "nonce", tx.Nonce())

	db.updateExecuteTxMetrics(tx)

	return db.ethState.ExecuteTx(tx)
}

// Persist finalises the current block and writes it to disk.
//
// Returns the persisted Block.
func (db *Database) Persist(receiver common.Address) (ethTypes.Block, error) {
	db.logger.Info("Persisting DB Block", "data", db.ethState.blockState)
	db.logger.Info(fmt.Sprintf("Current time in millisecond: %d", time.Now().UnixNano()/1000000))

	db.metrics.PersistedTxsTotal.Add(float64(len(db.ethState.blockState.transactions)))
	db.metrics.ChaindbHeight.Set(float64(db.ethState.blockState.header.Number.Uint64()))

	return db.ethState.Persist(receiver)
}

// ResetBlockState resets the intxBroadcastLoop-memory block's processing state.
func (db *Database) ResetBlockState(receiver common.Address) error {
	db.logger.Debug("Resetting DB BlockState", "receiver", receiver.Hex())
	return db.ethState.ResetBlockState(receiver)
}
type Header struct {
	// basic block info
	Time     time.Time `protobuf:"bytes,4,opt,name=time,stdtime" json:"time"`
	NumTxs   int64     `protobuf:"varint,5,opt,name=num_txs,json=numTxs,proto3" json:"num_txs,omitempty"`
	TotalTxs int64     `protobuf:"varint,6,opt,name=total_txs,json=totalTxs,proto3" json:"total_txs,omitempty"`

}
func (m *Header) GetNumTxs() int64 {
	if m != nil {
		return m.NumTxs
	}
	return 0
}
// UpdateBlockState uses the tendermint header to update the eth header.
func (db *Database) UpdateBlockState(tmHeader *Header) {
	db.logger.Debug("Updating DB BlockState")
	db.ethState.UpdateBlockState(
		db.eth.APIBackend.ChainConfig(),
		uint64(tmHeader.Time.Unix()),
		uint64(tmHeader.GetNumTxs()),
	)
}

// GasLimit returns the maximum gas per block.
func (db *Database) GasLimit() uint64 {
	return db.ethState.GasLimit().Gas()
}

// APIs returns the collection of Ethereum RPC services.
//
// Overwrites go-ethereum/eth/backend.go::APIs().
//
// Some of the API methods must be re-implemented to support Ethereum web3 features
// due to dependency on Tendermint, e.g Syncing().
func (db *Database) APIs() []rpc.API {
	ethAPIs := db.Ethereum().APIs()
	newAPIs := []rpc.API{}

	for _, v := range ethAPIs {
		if isDisabledAPI(v.Namespace) {
			continue
		}

		if _, ok := v.Service.(*eth.PublicMinerAPI); ok {
			continue
		}

		if v.Namespace == "net" {
			v.Service = dbAPI.NewPublicNetAPI(db.ethCfg.NetworkId)
		}

		if _, ok := v.Service.(*eth.PublicEthereumAPI); ok {
			v.Service = dbAPI.NewPublicEthereumAPI(
				db.ethCfg.Genesis.Config.ChainID,
				db.eth,
			)
		}

		if _, ok := v.Service.(downloader.PublicDownloaderAPI); ok {
			v.Service = dbAPI.NewPublicDownloaderAPI()
		}

		newAPIs = append(newAPIs, v)
	}

	return newAPIs
}

func (db *Database) Start(_ *p2p.Server) error {
	//go db.txBroadcastLoop()
	return nil
}

func (db *Database) Stop() error {
	//db.ethTxSub.Unsubscribe()
	db.eth.BlockChain().Stop()
	db.eth.Engine().Close()
	db.eth.TxPool().Stop()
	db.eth.ChainDb().Close()
	//if err := db.eth.Stop(); err != nil {
	//	return err
	//}
	return nil
}

func (db *Database) Protocols() []p2p.Protocol {
	return nil
}

func (db *Database) updateExecuteTxMetrics(tx *ethTypes.Transaction) {
	db.metrics.ExecutedTxsTotal.Add(1)

	txCost, _ := web3.WeiToPhoton(tx.Cost()).Float64()
	db.metrics.TxsCostTotal.Add(txCost)

	txGasWei := new(big.Int).Mul(big.NewInt(int64(tx.Gas())), tx.GasPrice())
	txGas, _ := web3.WeiToPhoton(txGasWei).Float64()
	db.metrics.TxsGasTotal.Add(txGas)

	db.metrics.TxsSizeTotal.Add(float64(tx.Size()))
}

func isDisabledAPI(namespace string) bool {
	return namespace == "miner" || namespace == "admin"
}

// Transactions sent via the go-ethereum rpc need to be routed to tendermint.
//
// Listening to txs and forward to tendermint.
//from ethraw directly route transactions here without adding to pool

//func SubmitFastTransaction(db *Database) {
	//ToApi(db)
	//go db.txBroadcastLoop()
//}

//use this also with another package

//func ToApi(mydb *Database){
	//txbypass.SenderApi(mydb)

//}


//perhaps move this out to another package to avoid the circular dependency
//func (db *Database) BypassPoolTx(tx *types.Transaction){
	//go db.TxloopAragorn(tx)
//}
//func (db *Database) TxloopAragorn(tx *types.Transaction){
	//ReceiveChan <- tx
//}
var (
	ProposeChan chan *types.Transaction
)

func (db *Database) TxBroadcastLoop() {

	db.ethTxsCh = make(chan core.NewTxsEvent, txChanSize)
	db.ethTxSub = db.eth.TxPool().SubscribeNewTxsEvent(db.ethTxsCh)
	//ProposeChan = make(chan *types.Transaction, 4096)
	for obj := range db.ethTxsCh {
		db.logger.Debug("Captured NewTxsEvent from pool")
		for _, tx := range obj.Txs {

			if db.start {
				fmt.Printf("START: %d\n", time.Now().UnixNano()/1000000)
				db.start = false
			}

			db.proposeChan <- tx
			//ProposeChan <- tx

		}
	}
}


func (db *Database) txBroadcastLoop() {
	db.ethTxsCh = make(chan core.NewTxsEvent, txChanSize)
	db.ethTxSub = db.eth.TxPool().SubscribeNewTxsEvent(db.ethTxsCh)

	for obj := range db.ethTxsCh {
		db.logger.Debug("Captured NewTxsEvent from pool")
		for _, tx := range obj.Txs {

			if db.start {
				fmt.Printf("START: %d\n", time.Now().UnixNano()/1000000)
				db.start = false
			}

			db.proposeChan <- tx

		}
	}
}
