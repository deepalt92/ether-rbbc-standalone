package node

import (
	//	"ether-rbbc/consensus"
	"ether-rbbc/database"
	"ether-rbbc/dbft"
	"ether-rbbc/log"
	"ether-rbbc/prometheus"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	tmtLog "github.com/ethereum/go-ethereum/log"
)

const (
	//TENDERMINT = iota
	RBBC = iota

	bufferSize = 256
)

type Node struct {
	cfg            *Config
	dbNode         *database.Node
//	consensusNode  *consensus.Node
	prometheusNode *prometheus.Node
	dbftNode       *dbft.Node
	logger         tmtLog.Logger
	conType        int
}

func NewNode(cfg *Config) (*Node, error) {
	conType := RBBC
	logger := log.NewLogger().With("engine", "node")

	logger.Info("Creating new prometheus node instance from config...")
	prometheusNode := prometheus.NewNode(cfg.prometheusCfg)

	//logger.Info("Creating new consensus node instance from config...")
	//consensusNode, err := consensus.NewNode(&cfg.consensusCfg, prometheusNode.Registry())
	/*if err != nil {
		return nil, err
	}*/

	proposeChan := make(chan *types.Transaction, bufferSize)
	// Todo: #90 Should be refactored, creating a new instance of a struct SHOULDN'T do any FS changes
	//consensusAPI := conAPI.NewConsensusApi(consensusNode.IsRunning)
	logger.Info("Creating new database node instance from config and creates a keystore dir...")
	dbNode, err := database.NewNode(&cfg.dbCfg, prometheusNode.Registry(), proposeChan)
	if err != nil {
		return nil, err
	}

	dbftNode := dbft.NewNode(0, proposeChan, int(cfg.BlockThreshold), int(cfg.BlockTimeout))

	return &Node{
		cfg:            cfg,
		dbNode:         dbNode,
		prometheusNode: prometheusNode,
		dbftNode:       dbftNode,
		logger:         logger,
		conType:        conType,
	}, nil
}

func (n *Node) Start(dbftdir string) error {
	n.logger.Info("Starting database engine...")
	if err := n.dbNode.Start(); err != nil {
		return err
	}

	n.logger.Info("Starting consensus engine...")

	switch n.conType {
	//case TENDERMINT:
	//	if err := n.consensusNode.Start(n.dbNode.RpcClient(), n.dbNode.Database()); err != nil {
	//		return err
	//	}
	case RBBC:
		n.logger.Info("Initialize DBFT Server")
		n.dbNode.Database().ResetBlockState(common.Address{})
		n.dbftNode.Start(n.dbNode.Database(), dbftdir)
		n.logger.Info("Started DBFT")
	}

	if err := n.prometheusNode.Start(); err != nil {
		return err
	}

	return nil
}

func (n *Node) Stop() error {
	// IMPORTANT: We need to close consensus first so that node stops receiving new blocks
	// before database is closed
	//n.logger.Info("Stopping consensus engine...")
	//if err := n.consensusNode.Stop(); err != nil {
		//return err
	//}
	//n.logger.Info("Consensus node stopped")

	n.logger.Info("Stopping database engine...")
	if err := n.dbNode.Stop(); err != nil {
		return err
	}
	n.logger.Info("Database node stopped")

	if err := n.prometheusNode.Stop(); err != nil {
		return err
	}

	return nil
}
