package database

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/trie"

	//tmtCode "github.com/tendermint/tendermint/abci/example/code"
	//tmtAbciTypes "github.com/tendermint/tendermint/abci/types"
	"math/big"
	"runtime"
)

// The blockState struct handles processing of TXs included in a block.
//
// It's updated with each ExecuteTx and reset on Persist.
type blockState struct {
	header *ethTypes.Header
	parent *ethTypes.Block
	state  *state.StateDB

	txIndex      int
	transactions []*ethTypes.Transaction
	receipts     ethTypes.Receipts

	totalUsedGas uint64
	gp           *core.GasPool
}

// Executes TX against the eth blockchain state.
//
// Fetches TX logs and returns new TX receipt. The changes happen only
// inside of the Eth state, not disk!
//
// Logic copied from `core/state_processor.go` `(p *StateProcessor) Process` that gets
// normally executed on block persist.
func (bs *blockState) execTx(bc *core.BlockChain, config *eth.Config, chainConfig *params.ChainConfig, blockHash common.Hash, tx *ethTypes.Transaction) error {
	// TODO: Investigate if snapshot should be used `snapshot := bs.state.Snapshot()`
	bs.state.Prepare(tx.Hash(), bs.txIndex)
	fmt.Println("Trying to apply transaction")

	receipt, err := core.ApplyTransaction(
		chainConfig,
		bc,
		nil, // defaults to address of the author of the header
		bs.gp,
		bs.state,
		bs.header,
		tx,
		&bs.totalUsedGas,
		vm.Config{EnablePreimageRecording: config.EnablePreimageRecording},
	)
	if err != nil {
		fmt.Println("Error applying transaction", err)
		// TODO: investigate if snapshot should be used `bs.state.RevertToSnapshot(snapshot)`
		//return tmtAbciTypes.ResponseDeliverTx{Code: tmtCode.CodeTypeEncodingError, Log: fmt.Sprintf("Error applying state TX %v", err)}
		return err
	}

	bs.txIndex++
	//fmt.Println("The receipt status is:", receipt.Status)
	// The slices are allocated in updateBlockState
	bs.transactions = append(bs.transactions, tx)
	bs.receipts = append(bs.receipts, receipt)

	//PrintMemUsage()
	//fmt.Println("Total system memory: %d\n", memory.TotalMemory())
	//return tmtAbciTypes.ResponseDeliverTx{Code: tmtAbciTypes.CodeTypeOK}
	return nil
}
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys)) //in mega bytes
	fmt.Printf("\tNumGC = %v\n", m.NumGC) // number of garbage collection cycles
}
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

// Persist the eth sate, update the header, make a new block and save it to disk.
//
// Returns the persisted Block.

func (bs *blockState) persist(bc *core.BlockChain, db ethdb.Database) (ethTypes.Block, error) {


	rootHash, err := bs.state.Commit(false)
	if err != nil {
		return ethTypes.Block{}, err
	}
	bs.header.Root = rootHash

	// Write block to disk
	//you can perhaps include the hash of the transactions here -- this is actually in block.go//
	block := ethTypes.NewBlock(bs.header, bs.transactions, nil, bs.receipts, trie.NewStackTrie(nil))

	//caching block body
	//rawdb.BlockBody[bs.header.Number.Uint64()]=block.Body()
	rawdb.BlockBody.Store(bs.header.Number.Uint64(), block.Body())
	//caching receipts
	core.CacheReceipts(bc, block, bs.receipts)
	_, err = bc.MyInsertChain([]*ethTypes.Block{block}, bs.receipts)
	if err != nil {
		return ethTypes.Block{}, err
	}

	return *block, nil
}

func (bs *blockState) updateBlockState(config *params.ChainConfig, parentTime uint64, numTx uint64) {
	//parentHeader := bs.parent.Header()
	bs.header.Time = new(big.Int).SetUint64(parentTime).Uint64()
	//bs.header.Difficulty = ethash.CalcDifficulty(config, parentTime, parentHeader)
	bs.header.Difficulty = big.NewInt(1)
	// you can reset them instead of doing make again!!!!
		bs.transactions = make([]*ethTypes.Transaction, 0, numTx)
		bs.receipts = make([]*ethTypes.Receipt, 0, numTx)

}
