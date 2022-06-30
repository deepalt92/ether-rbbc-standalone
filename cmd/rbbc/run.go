package main

import (
	"ether-rbbc/database"
	"ether-rbbc/log"
	"ether-rbbc/node"
	"ether-rbbc/prometheus"
	"ether-rbbc/tracer"
	"fmt"
	ethLog "github.com/ethereum/go-ethereum/log"
	"github.com/spf13/cobra"
  //   "github.com/tendermint/tendermint/libs/common"
	"os/signal"
	"syscall"

	//"github.com/tendermint/tendermint/libs/common"
	"gopkg.in/urfave/cli.v1"
	"os"
	"path"
	"path/filepath"
	//"rbbc/configurations"
)

const (
	//TendermintP2PListenPort = uint(26656)
	//TendermintRpcListenPort = uint(26657)
	//TendermintProxyAppName  = "rbbc"
	DBFTDefaultDir          = ".rbbc/config"
)
var dbftConfigDir string
var (
	/*ConsensusRpcListenPortFlag = cli.UintFlag{
		Name:  "tmt_rpc_port",
		Value: TendermintRpcListenPort,
		Usage: "Tendermint RPC port used to receive incoming messages from Lightchain",
	}
	ConsensusP2PListenPortFlag = cli.UintFlag{
		Name:  "tmt_p2p_port",
		Value: TendermintP2PListenPort,
		Usage: "Tendermint port used to achieve exchange messages across nodes",
	}
	ConsensusProxyAppNameFlag = cli.StringFlag{
		Name:  "abci_name",
		Value: TendermintProxyAppName,
		Usage: "socket | grpc",
	}*/
	PrometheusFlag = cli.BoolFlag{
		Name:  "prometheus",
		Usage: "Enable prometheus metrics exporter",
	}
	DbftFlag = cli.StringFlag{
		Name:  "dbftDir",
		Usage: "Configuration file for DBFT consensus system",
		Value: fmt.Sprintf("%s/%s", defaultHomeDir, DBFTDefaultDir),
	}
	BlockThreshold = cli.UintFlag{
		Name:  "threshold",
		Usage: "Threshold blocksize",
		Value: 100,
	}
	BlockTimeout = cli.UintFlag{
		Name:  "timeout",
		Usage: "Threshold blocksize",
		Value: 1000,
	}
)
func TrapSignal(logger ethLog.Logger, cb func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for sig := range c {
			logger.Info(fmt.Sprintf("captured %v, exiting...", sig))
			if cb != nil {
				cb()
			}
			os.Exit(0)
		}
	}()
}

func runCmd() *cobra.Command {
	var runCmd = &cobra.Command{
		Use:   "run",
		Short: "Launches rbbc node and all of its online services including blockchain (Geth) and DBFT consensus.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			lvlStr, _ := cmd.Flags().GetString(LogLvlFlag.Name)
			if lvl, err := ethLog.LvlFromString(lvlStr); err == nil {
				log.SetupLogger(lvl)
			}

			logger.Info("Launching Lightchain node...")

			dataDir, _ := cmd.Flags().GetString(DataDirFlag.GetName())
			shouldTrace, _ := cmd.Flags().GetBool(TraceFlag.Name)
			//rpcListenPort, _ := cmd.Flags().GetUint(ConsensusRpcListenPortFlag.GetName())
			//p2pListenPort, _ := cmd.Flags().GetUint(ConsensusP2PListenPortFlag.GetName())
			//proxyAppName, _ := cmd.Flags().GetString(ConsensusProxyAppNameFlag.GetName())
			enablePrometheus, _ := cmd.Flags().GetBool(PrometheusFlag.GetName())
			databaseDataDir := filepath.Join(dataDir, database.DataDirPath)
			dbftConfigDir, _ = cmd.Flags().GetString(DbftFlag.GetName())
			threshold, _ := cmd.Flags().GetUint(BlockThreshold.GetName())
			timeout, _ := cmd.Flags().GetUint(BlockTimeout.GetName())

			/*consensusCfg, err := consensus.NewConfig(
				filepath.Join(dataDir, consensus.DataDirName),
				rpcListenPort,
				p2pListenPort,
				proxyAppName,
				enablePrometheus,
			)*/
			/*if err != nil {
				logger.Error(fmt.Errorf("consensus node config could not be created: %v", err).Error())
				os.Exit(1)
			}*/

			// Fake cli.context required by Ethereum node
			ctx := newNodeClientCtx(databaseDataDir, cmd)
			dbCfg, err := database.NewConfig(databaseDataDir, enablePrometheus, ctx)
			if err != nil {
				logger.Error(fmt.Errorf("database node config could not be created: %v", err).Error())
				os.Exit(1)
			}

			prometheusCfg := prometheus.NewConfig(
				enablePrometheus,
				prometheus.DefaultPrometheusAddr,
				prometheus.DefaultPrometheusNamespace,
				dbCfg.GethIpcPath(),
			)

			tracerCfg := tracer.NewConfig(shouldTrace, path.Join(dataDir, "tracer.log"))

			if shouldTrace {
				tracerCfg.PrintWarning(logger)
			}

			//observerConfig := new(configurations.ObserverConfig)
			//configurations.ReadConfig(fmt.Sprintf("%s/observer.yaml", dbftConfigDir), observerConfig)

			nodeCfg := node.NewConfig(dataDir, dbCfg, prometheusCfg, tracerCfg, threshold, timeout)

			n, err := node.NewNode(&nodeCfg)
			//if the configuration files change, change the node.//
			if err != nil {
				logger.Error(fmt.Errorf("rbbc node could not be instantiated: %v", err).Error())
				os.Exit(1)
			}

			TrapSignal(logger, func() {
				if err := n.Stop(); err != nil {
					logger.Error(fmt.Errorf("error stopping rbbc node. %v", err).Error())
					os.Exit(1)
				}

				os.Exit(0)
			})

			logger.Debug("Starting rbbc node...")
			if err := n.Start(dbftConfigDir); err != nil {
				logger.Error(fmt.Errorf("rbbc node could not be started: %v", err).Error())
				os.Exit(1)
			}

			select {}
		},
	}

	addRunCmdFlags(runCmd)

	return runCmd
}


func addRunCmdFlags(cmd *cobra.Command) {
	addDefaultFlags(cmd)
	addConsensusFlags(cmd)
	addEthNodeFlags(cmd)
	cmd.Flags().String(DbftFlag.GetName(), DbftFlag.Value, DbftFlag.Usage)
	cmd.Flags().Uint(BlockThreshold.GetName(), BlockThreshold.Value, BlockThreshold.Usage)
	cmd.Flags().Uint(BlockTimeout.GetName(), BlockTimeout.Value, BlockTimeout.Usage)
	cmd.Flags().Bool(PrometheusFlag.GetName(), false, PrometheusFlag.Usage)
}

func addConsensusFlags(cmd *cobra.Command) {
	//cmd.Flags().Uint(ConsensusRpcListenPortFlag.GetName(), ConsensusRpcListenPortFlag.Value, ConsensusRpcListenPortFlag.Usage)
	//cmd.Flags().Uint(ConsensusP2PListenPortFlag.GetName(), ConsensusP2PListenPortFlag.Value, ConsensusP2PListenPortFlag.Usage)
}
