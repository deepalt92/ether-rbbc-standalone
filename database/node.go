package database

import (
	"ether-rbbc/database/metrics"
	"ether-rbbc/log"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	tmtLog "github.com/ethereum/go-ethereum/log"
	ethNode "github.com/ethereum/go-ethereum/node"
	ethRpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/prometheus/client_golang/prometheus"
	"os"
)

// Node is the main object.
type Node struct {
	ethereum  *ethNode.Node
	database  *Database
	rpcClient *ethRpc.Client
	cfg       *Config
	logger    tmtLog.Logger
}

// NewNode creates a new node.
func NewNode(cfg *Config, registry *prometheus.Registry, proposeChan chan *types.Transaction) (*Node, error) {
	logger := log.NewLogger().With("engine", "database")

	// @TODO Investigate why Genesis file is not automatically loaded
	var err error
	cfg.GethCfg.EthCfg.Genesis, err = readGenesisFile(cfg.genesisPath())
	if err != nil {
		return nil, err
	}

	// Todo: #90 Should be refactored, creating a new instance of a struct SHOULDN'T do any FS changes
	cfg.GethCfg.EthCfg.NetworkId = cfg.GethCfg.EthCfg.Genesis.Config.ChainID.Uint64()
	ethereum, err := ethNode.New(&cfg.GethCfg.NodeCfg)
	if err != nil {
		return nil, err
	}

	var trackedMetrics metrics.Metrics
	if cfg.metrics {
		trackedMetrics = metrics.NewMetrics(registry)
	} else {
		trackedMetrics = metrics.NewNullMetrics()
	}

	n := Node{
		ethereum,
		nil,
		nil,
		cfg,
		logger,
	}

	logger.Debug("Binding ethereum events to rpc client...")

	//stack.RegisterAPIs(tracers.APIs(backend.APIBackend))

	/*err = ethereum.Register(func(ctx *ethNode.ServiceContext) (ethNode.Service, error) {
		logger.Debug(fmt.Sprintf("Registering database..."))

		n.database, err = NewDatabase(ctx, &cfg.GethCfg.EthCfg, logger, trackedMetrics, proposeChan)
		if err != nil {
			return nil, err
		}

		return n.database, nil
	})
	if err != nil {
		return nil, err
	}*/
	/*datadir := "/home/deepal/.lightchain_standalone/"
	//datadir := os.MkdirAll("", "")
	config := &ethNode.Config{
		Name:    "geth",
		Version: params.Version,
		DataDir: datadir,
		P2P: p2p.Config{
			ListenAddr:  "0.0.0.0:0",
			NoDiscovery: true,
			MaxPeers:    25,
		},
		UseLightweightKDF: true,
	}
	// Create the node and configure a full Ethereum node on it
	stack, err := ethNode.New(config)*/
	n.database, err = NewDatabase(ethereum, &cfg.GethCfg.EthCfg, logger, trackedMetrics, proposeChan)

	if err != nil {
		return nil, err
	}

	return &n, nil
}

// Start starts base node and stops p2p server
func (n *Node) Start() error {
	n.logger.Info("Starting ethereum node")
	err := n.ethereum.Start()
	if err != nil {
		return err
	}

	n.logger.Info("Stopping p2p communication of ethereum node")
	n.ethereum.Server().Stop()

	n.logger.Info("Register wallet event handlers to open and auto-derive wallets")
	events := make(chan accounts.WalletEvent, 16)
	n.ethereum.AccountManager().Subscribe(events)
	go registerEventHandlers(n, events)

	// Fetch the registered service of this type
	//var ethDb *Database
	/*n.logger.Info("Validate ethereum service is running...")
	if err := n.ethereum.Service(&ethDb); err != nil {
		n.logger.Error(fmt.Errorf("ethereum service not running: %v", err).Error())
		return err
	}*/

	n.rpcClient, err = n.ethereum.Attach()
	if err != nil {
		n.logger.Error(fmt.Errorf("failed to attach to the inproc geth: %v", err).Error())
		os.Exit(1)
	}

	n.logger.Info("Started database engine...")
	return nil
}

func (n *Node) Stop() error {
	n.logger.Debug("Stopping ethereum service")
	/*if err := n.ethereum.Stop(); err != nil {
		return err
	}*/
	if err := n.ethereum.Close(); err != nil {
		return err
	}
	return nil
}

func (n *Node) RpcClient() *ethRpc.Client {
	return n.rpcClient
}

func (n *Node) Database() *Database {
	return n.database
}

func registerEventHandlers(n *Node, events chan accounts.WalletEvent) error {
	rpcClient, err := n.ethereum.Attach()
	if err != nil {
		n.logger.Error(fmt.Errorf("failed to attach to self: %v", err).Error())
		return err
	}
	stateReader := ethclient.NewClient(rpcClient)

	for _, wallet := range n.ethereum.AccountManager().Wallets() {
		if err := wallet.Open(""); err != nil {
			n.logger.Error("Failed to open wallet", "url", wallet.URL(), "err", err)
		} else {
			wallet.SelfDerive([]accounts.DerivationPath{accounts.DefaultBaseDerivationPath}, stateReader)
		}
	}

	// Listen for wallet event till termination
	for event := range events {
		switch event.Kind {
		case accounts.WalletArrived:
			if err := event.Wallet.Open(""); err != nil {
				n.logger.Error("New wallet appeared, failed to open", "url", event.Wallet.URL(), "err", err)
			}
		case accounts.WalletOpened:
			status, _ := event.Wallet.Status()
			n.logger.Info("New wallet appeared", "url", event.Wallet.URL(), "status", status)

			if event.Wallet.URL().Scheme == "ledger" {
				event.Wallet.SelfDerive([]accounts.DerivationPath{accounts.LegacyLedgerBaseDerivationPath}, stateReader)
			} else {
				event.Wallet.SelfDerive([]accounts.DerivationPath{accounts.DefaultBaseDerivationPath}, stateReader)
			}

		case accounts.WalletDropped:
			n.logger.Info("Old wallet dropped", "url", event.Wallet.URL())
			event.Wallet.Close()
		}
	}

	return nil
}
