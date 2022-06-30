package dbft

import (
	"ether-rbbc/database"
	"ether-rbbc/log"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	tmtLog "github.com/ethereum/go-ethereum/log"
//	N "rbbc/node"
	"time"
)

const (
	timerBufferSize  = 128
	commitBufferSize = 128
	txChanSize = 4096
)
//var super = make([]*types.Transaction, 100000)
type Node struct {
	id               int
	started          bool
	addr             string
	proposeChan      chan *types.Transaction
	commitChan       chan [][]byte
//	observerServer   *grpc.ObserverServer
//	consensusClient  *grpc.ConsensusClient
	db               *database.Database
	txm              *txManager
	logger           tmtLog.Logger
	threshold        int
	timeout          int
	timeoutChan      chan bool
	restartTimerChan chan bool
}

func NewNode(
	id int,
	proposeChan chan *types.Transaction,
	threshold int,
	timeout int) *Node {
	commitChan := make(chan [][]byte, commitBufferSize)
	return &Node{
		id:              id,
//		addr:            fmt.Sprintf("%s:%d", observerConfig.Ip, observerConfig.Port),
		commitChan:      commitChan,
//		observerServer:  grpc.NewObserverServer(observerConfig.GenerateServerConfig(), commitChan),
//		consensusClient: grpc.NewConsensusClient(observerConfig.GenerateClientConfig()),
		txm: &txManager{
			threshold: threshold,
		},
		logger:           log.NewLogger().With("engine", "observer"),
		proposeChan:      proposeChan,
		threshold:        threshold,
		timeout:          timeout,
		timeoutChan:      make(chan bool, timerBufferSize),
		restartTimerChan: make(chan bool, timerBufferSize),
	}
}


var (
	total = 0
	count = 0
	props [][]byte
	MyNumTx = 0
	execTimeProposal = int64(0)
	CommitChan = make(chan [][]byte, 4096)
)
func (node *Node) Start(db *database.Database, dbftdir string) {
	node.db = db

	if node.started {
		return
	}

	node.logger.Info(fmt.Sprintf("Starting DBFT Server. Block Threshold: %d, Block Timeout: %d", node.threshold, node.timeout))

	//node.observerServer.Start()

	//node.consensusClient.Subscribe(node.addr)

	//dbftConfig := new(configurations.DbftConfig)
	//configurations.ReadConfig(fmt.Sprintf("%s/dbft.yaml", dbftdir), dbftConfig)


	/*if dbftConfig.Port > 0 && dbftConfig.GrpcPort < 65536 {
		dbftConfig.GrpcPort = consensusPort
	}*/

	/*if dbftConfig.Id >= 0 && dbftConfig.Id < dbftConfig.Size {
		//dbftConfig.Id = dbftConfig.Id
		dbftConfig.Port = 1888 + dbftConfig.Id*1000
	}*/

	//rbbc := N.NewNode(dbftConfig)

	//perf.StartReporting(5 * time.Second)

	/*go func() {
		rbbc.Start()
	}()*/



	//ethapi.ReceiveChan = make(chan *types.Transaction, 4096)


	// Listen to commitChan for decided super blocks
	// and commit super block upon reception.
	go func() {

		//for block := range node.commitChan {
		for block := range CommitChan {
			//tot := 0
			total := 0

			fmt.Printf("COMMIT: %d\n", time.Now().UnixNano()/1000000)

			for _, proposal := range block {
				//super = append(super,proposal...)
				txs := node.txm.deserialize(proposal)
				//super = append(super, txs...)
				node.db.UpdateBlockState(&database.Header{Time: time.Now(), NumTxs: int64(len(txs))})
                //startT := time.Now().UnixNano()/1000000
				for _, tx := range txs {
					node.db.ExecuteTx(tx)
				}
                //endT:= time.Now().UnixNano()/1000000
                //execTimeProposal = execTimeProposal+ endT - startT
               // fmt.Println("Proposal execution time:", execTimeProposal)
				node.db.Persist(common.Address{})
				total = len(txs)
			}
			fmt.Printf("Number of Transactions: %d\n", total)
			fmt.Println("Proposal execution time:", execTimeProposal)
		}
	}()
	// Timer go routine tracks timeout status.
		// A timer waits for certain timeout and notifie}s timeoutChan for timeout event.
	// Timer is reset when timer received message from restartTimerChan/
	go func() {
		for {
			select {
			case <-time.After(time.Duration(node.timeout) * time.Millisecond):
				node.timeoutChan <- true
			case <-node.restartTimerChan:
			}
		}
	}()

	go func() {
		for {
			select {
			case tx := <- core.ProposeChan:
			//case tx := <-node.proposeChan:
				MyNumTx++
				fmt.Println("received transactions at propose channel", MyNumTx)
				node.txm.add(tx)
				if node.txm.check() {
					// Reset timer after a proposal is made
					go func() {
						node.restartTimerChan <- true
					}()

					// Make proposal
					//data := node.txm.serialize()
					//dbft.PChannel <- data
                    /////////////////////////////////////////////////

					props = append(props, node.txm.serialize())
					CommitChan <- props
					props = nil

				}
			case <-node.timeoutChan:
				// Make proposal is timer has expired
				// and tx pool is not empty.
				if node.txm.size() > 0 {
					//node.consensusClient.Propose(node.txm.serialize())
					///////////////////////////////////
					//data := node.txm.serialize()
					//dbft.PChannel <- data
					//////////////////////////////

					props = append(props, node.txm.serialize())
					CommitChan <- props
					props = nil

				}
			}
		}
	}()

}
