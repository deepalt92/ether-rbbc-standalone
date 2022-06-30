package main

import (
	"ether-rbbc/log"
	"fmt"
	ethUtils "github.com/ethereum/go-ethereum/cmd/utils"
	ethLog "github.com/ethereum/go-ethereum/log"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"gopkg.in/urfave/cli.v1"
	"net/http"
	"os"
	"path"
	"runtime/debug"
    _ "net/http/pprof"
)

var logger = log.NewLogger().With("module", "rbbc")

var (
	defaultHomeDir, _ = homedir.Dir()

	DataDirFlag = ethUtils.DirectoryFlag{
		Name:  "datadir",
		Usage: "Data directory for the databases and keystore",
		Value: ethUtils.DirectoryString(path.Join(defaultHomeDir, ".lightchain")),
	}

	LogLvlFlag = cli.StringFlag{
		Name:  "lvl",
		Usage: "Level of logging",
		Value: ethLog.LvlInfo.String(),
	}

	TraceFlag = cli.BoolFlag{
		Name:  "trace",
		Usage: "Whenever to be asserting and reporting blockchain state in real-time (testing, debugging purposes)",
	}
)

func main() {

	lightchainCmd := LightchainCmd()

	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Node resulted in panic: %s. \n"+string(debug.Stack()), r)
			os.Exit(1)
		}
	}()

	if err := lightchainCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// LightchainCmd is the main Lightstreams PoA blockchain node.
func LightchainCmd() *cobra.Command {
	var lightchainCmd = &cobra.Command{
		Use:   "rbbc",
		Short: "Lightstreams PoA blockchain node.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return incorrectUsageErr()
		},
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	lightchainCmd.AddCommand(versionCmd)
	lightchainCmd.AddCommand(docsCmd())
	lightchainCmd.AddCommand(initCmd())
	lightchainCmd.AddCommand(runCmd())
	lightchainCmd.AddCommand(simulateCmd())

	return lightchainCmd
}

func addDefaultFlags(cmd *cobra.Command) {
	cmd.Flags().String(DataDirFlag.Name, DataDirFlag.Value.String(), DataDirFlag.Usage)
	cmd.Flags().String(LogLvlFlag.Name, LogLvlFlag.Value, LogLvlFlag.Usage)

	cmd.Flags().Bool(TraceFlag.Name, false, TraceFlag.Usage)
}

func incorrectUsageErr() error {
	return fmt.Errorf("incorrect usage. More instructions also available at https://docs.lightstreams.network/cli-docs/rbbc/")
}
