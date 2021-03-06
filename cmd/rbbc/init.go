package main

import (
	"ether-rbbc/database"
	"ether-rbbc/fs"
	"ether-rbbc/log"
	"ether-rbbc/network"
	"ether-rbbc/node"
	"ether-rbbc/prometheus"
	"ether-rbbc/tracer"
	"fmt"
	ethLog "github.com/ethereum/go-ethereum/log"
	"github.com/spf13/cobra"
	"gopkg.in/urfave/cli.v1"
	"os"
	"path"
	"path/filepath"
)

var (
	StandAloneNetFlag = cli.BoolFlag{
		Name:  "standalone",
		Usage: "Initialize a stand alone node not connected to any network",
	}

	SiriusNetFlag = cli.BoolFlag{
		Name:  "sirius",
		Usage: "Initialize a node connected to Sirius network",
	}

	MainNetFlag = cli.BoolFlag{
		Name:  "mainnet",
		Usage: "Initialize a node connected to production, MainNet",
	}

	ForceFlag = cli.BoolFlag{
		Name:  "force",
		Usage: "Forces the init by removing the data dir if already exists",
	}
)

func initCmd() *cobra.Command {
	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initializes new rbbc node according to the configured flags.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: initCmdRun,
	}

	addDefaultFlags(initCmd)
	addInitCmdFLags(initCmd)

	return initCmd
}

func initCmdRun(cmd *cobra.Command, args []string) {
	nodeCfg, ntw, err := newNodeCfgFromCmd(cmd)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	if err := node.Init(nodeCfg, ntw); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.Info(fmt.Sprintf("Lightchain node successfully initialized into '%s'!", nodeCfg.DataDir))
	os.Exit(0)
}

func newNodeCfgFromCmd(cmd *cobra.Command) (node.Config, network.Network, error) {
	lvlStr, _ := cmd.Flags().GetString(LogLvlFlag.Name)
	if lvl, err := ethLog.LvlFromString(lvlStr); err == nil {
		log.SetupLogger(lvl)
	}

	dataDir, _ := cmd.Flags().GetString(DataDirFlag.Name)
	forceInit, _ := cmd.Flags().GetBool(ForceFlag.Name)
	shouldTrace, _ := cmd.Flags().GetBool(TraceFlag.Name)

	// To prevent accidental data dir like "/" and disasters
	if len(dataDir) < 3 {
		return node.Config{}, "", fmt.Errorf("DataDir must be absolute path of at least 3 chars")
	}

	// This should be done inside of the `node.Init()` pkg but due to bad design,
	// creation of new Ethereum node instance from a config accidentally modifies the FS by creating
	// a keystore dir therefore we have to perform this check as early as possible.
	//
	// Todo: After #90 is done, move this code to `Node.Init()`
	if isEmpty, err := fs.IsDirEmptyOrNotExists(dataDir); !isEmpty || err != nil {
		if err != nil {
			return node.Config{}, "", fmt.Errorf("unable to initialize rbbc node. %s", err.Error())
		}

		if !forceInit {
			return node.Config{}, "", fmt.Errorf("unable to initialize rbbc node. %s already exists", dataDir)
		}

		isConfirmed := fs.AskForConfirmation(fmt.Sprintf("Are you sure to erase %s ?", dataDir))
		if !isConfirmed {
			return node.Config{}, "", fmt.Errorf("unable to initialize rbbc node. %s already exists", dataDir)
		}

		logger.Info(fmt.Sprintf("Removing '%s' folder ...", dataDir))
		if err := fs.RemoveAll(dataDir); err != nil {
			return node.Config{}, "", fmt.Errorf("unable to remove data dir '%s'. %s", dataDir, err)
		}
	}

	ntw, err := chooseNetwork(cmd)
	if err != nil {
		return node.Config{}, "", err
	}

	/*consensusCfg, err := consensus.NewConfig(
		filepath.Join(dataDir, consensus.DataDirName),
		TendermintRpcListenPort,
		TendermintP2PListenPort,
		TendermintProxyAppName,
		false,
	)*/

	if err != nil {
		return node.Config{}, "", err
	}

	dbDataDir := filepath.Join(dataDir, database.DataDirPath)
	dbCfg, err := database.NewConfig(dbDataDir, false, newNodeClientCtx(dbDataDir, cmd))

	if err != nil {
		return node.Config{}, "", err
	}

	prometheusCfg := prometheus.NewConfig(
		false,
		prometheus.DefaultPrometheusAddr,
		prometheus.DefaultPrometheusNamespace,
		dbCfg.GethIpcPath(),
	)

	tracerCfg := tracer.NewConfig(shouldTrace, path.Join(dataDir, "tracer.log"))

	if shouldTrace {
		tracerCfg.PrintWarning(logger)
	}

	return node.NewConfig(dataDir, dbCfg, prometheusCfg, tracerCfg, 0, 0), ntw, nil
}

func chooseNetwork(cmd *cobra.Command) (network.Network, error) {
	useStandAloneNet, _ := cmd.Flags().GetBool(StandAloneNetFlag.Name)
	useSiriusNet, _ := cmd.Flags().GetBool(SiriusNetFlag.Name)
	useMainNet, _ := cmd.Flags().GetBool(MainNetFlag.Name)

	if (boolToInt(useStandAloneNet) + boolToInt(useSiriusNet) + boolToInt(useMainNet)) > 1 {
		return "", fmt.Errorf("select only one network")
	}

	//if useStandAloneNet {
		return network.StandaloneNetwork, nil
	//}

	/*if useSiriusNet {
		return network.SiriusNetwork, nil
	}

	return network.MainNetNetwork, nil*/
}

func boolToInt(a bool) int {
	if a {
		return 1
	}

	return 0
}

func addInitCmdFLags(cmd *cobra.Command) {
	cmd.Flags().Bool(StandAloneNetFlag.Name, false, DataDirFlag.Usage)
	cmd.Flags().Bool(SiriusNetFlag.Name, false, SiriusNetFlag.Usage)
	cmd.Flags().Bool(MainNetFlag.Name, false, MainNetFlag.Usage)
	cmd.Flags().Bool(ForceFlag.Name, false, ForceFlag.Usage)
}
