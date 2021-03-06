package main

import (
	"flag"
	ethUtils "github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/spf13/cobra"
	"gopkg.in/urfave/cli.v1"
	"strconv"
)

var (
	RPCEnabledFlag       = ethUtils.HTTPEnabledFlag
	RPCListenAddrFlag    = ethUtils.HTTPListenAddrFlag
	RPCPortFlag          = ethUtils.HTTPPortFlag
	RPCApiFlag           = ethUtils.HTTPApiFlag
	WSEnabledFlag        = ethUtils.WSEnabledFlag
	WSListenAddrFlag     = ethUtils.WSListenAddrFlag
	WSPortFlag           = ethUtils.WSPortFlag
	WSApiFlag            = ethUtils.WSApiFlag
	RPCCORSDomainFlag    = ethUtils.HTTPCORSDomainFlag
	RPCVirtualHostsFlag  = ethUtils.HTTPVirtualHostsFlag
	WSAllowedOriginsFlag = ethUtils.WSAllowedOriginsFlag
)

func addEthNodeFlags(cmd *cobra.Command) {
	// RPC Flags
	cmd.Flags().Bool(RPCEnabledFlag.GetName(), false, RPCEnabledFlag.Usage)
	cmd.Flags().String(RPCListenAddrFlag.GetName(), RPCListenAddrFlag.Value, RPCListenAddrFlag.Usage)
	cmd.Flags().Int(RPCPortFlag.GetName(), RPCPortFlag.Value, RPCPortFlag.Usage)
	cmd.Flags().String(RPCApiFlag.GetName(), RPCApiFlag.Value, RPCApiFlag.Usage)
	cmd.Flags().String(RPCCORSDomainFlag.GetName(), RPCCORSDomainFlag.Value, RPCCORSDomainFlag.Usage)
	cmd.Flags().String(RPCVirtualHostsFlag.GetName(), RPCVirtualHostsFlag.Value, RPCVirtualHostsFlag.Usage)

	// WS Flags
	cmd.Flags().Bool(WSEnabledFlag.GetName(), false, WSEnabledFlag.Usage)
	cmd.Flags().String(WSListenAddrFlag.GetName(), WSListenAddrFlag.Value, WSListenAddrFlag.Usage)
	cmd.Flags().Int(WSPortFlag.GetName(), WSPortFlag.Value, WSPortFlag.Usage)
	cmd.Flags().String(WSApiFlag.GetName(), WSApiFlag.Value, WSApiFlag.Usage)
	cmd.Flags().String(WSAllowedOriginsFlag.GetName(), WSAllowedOriginsFlag.Usage, WSAllowedOriginsFlag.Usage)
}

func newNodeClientCtx(dataDir string, cmd *cobra.Command) *cli.Context {

	app := ethUtils.NewApp("0.0.0", "the rbbc command line interface")
	flagSet := flag.NewFlagSet("Ethereum CLI ctx", flag.ExitOnError)

	flagSet.String(ethUtils.DataDirFlag.GetName(), dataDir, ethUtils.DataDirFlag.Usage)

	rpcEnabledFlagValue, _ := cmd.Flags().GetBool(RPCEnabledFlag.GetName())
	flagSet.Bool(RPCEnabledFlag.GetName(), rpcEnabledFlagValue, RPCEnabledFlag.Usage)
	flagSet.Set(RPCEnabledFlag.GetName(), strconv.FormatBool(rpcEnabledFlagValue))

	rpcListenAddrValue, _ := cmd.Flags().GetString(RPCListenAddrFlag.GetName())
	flagSet.String(RPCListenAddrFlag.GetName(), rpcListenAddrValue, RPCListenAddrFlag.Usage)
	flagSet.Set(RPCListenAddrFlag.GetName(), rpcListenAddrValue)

	rpcPortFlagValue, _ := cmd.Flags().GetInt(RPCPortFlag.GetName())
	flagSet.Int(RPCPortFlag.GetName(), rpcPortFlagValue, RPCPortFlag.Usage)
	flagSet.Set(RPCPortFlag.GetName(), strconv.Itoa(rpcPortFlagValue))

	rpcApiFlagValue, _ := cmd.Flags().GetString(RPCApiFlag.GetName())
	flagSet.String(RPCApiFlag.GetName(), rpcApiFlagValue, RPCApiFlag.Usage)
	flagSet.Set(RPCApiFlag.GetName(), rpcApiFlagValue)

	wsEnabledFlagValue, _ := cmd.Flags().GetBool(WSEnabledFlag.GetName())
	flagSet.Bool(WSEnabledFlag.GetName(), wsEnabledFlagValue, WSEnabledFlag.Usage)
	flagSet.Set(WSEnabledFlag.GetName(), strconv.FormatBool(wsEnabledFlagValue))

	wsListenAddrFlagValue, _ := cmd.Flags().GetString(WSListenAddrFlag.GetName())
	flagSet.String(WSListenAddrFlag.GetName(), wsListenAddrFlagValue, WSListenAddrFlag.Usage)
	flagSet.Set(WSListenAddrFlag.GetName(), wsListenAddrFlagValue)

	wsPortFlagValue, _ := cmd.Flags().GetInt(WSPortFlag.GetName())
	flagSet.Int(WSPortFlag.GetName(), wsPortFlagValue, WSPortFlag.Usage)
	flagSet.Set(WSPortFlag.GetName(), strconv.Itoa(wsPortFlagValue))

	wsApiFlagValue, _ := cmd.Flags().GetString(WSApiFlag.GetName())
	flagSet.String(WSApiFlag.GetName(), wsApiFlagValue, WSApiFlag.Usage)
	flagSet.Set(WSApiFlag.GetName(), wsApiFlagValue)

	rpcCORSDomainFlagValue, _ := cmd.Flags().GetString(RPCCORSDomainFlag.GetName())
	flagSet.String(RPCCORSDomainFlag.GetName(), rpcCORSDomainFlagValue, RPCCORSDomainFlag.Usage)
	flagSet.Set(RPCCORSDomainFlag.GetName(), rpcCORSDomainFlagValue)

	rpcVirtualHostsFlagValue, _ := cmd.Flags().GetString(RPCVirtualHostsFlag.GetName())
	flagSet.String(RPCVirtualHostsFlag.GetName(), rpcVirtualHostsFlagValue, RPCVirtualHostsFlag.Usage)
	flagSet.Set(RPCVirtualHostsFlag.GetName(), rpcVirtualHostsFlagValue)

	wsAllowedOriginsFlagValue, _ := cmd.Flags().GetString(WSAllowedOriginsFlag.GetName())
	flagSet.String(WSAllowedOriginsFlag.GetName(), wsAllowedOriginsFlagValue, WSAllowedOriginsFlag.Usage)
	flagSet.Set(WSAllowedOriginsFlag.GetName(), wsAllowedOriginsFlagValue)

	// Default values required
	flagSet.String(ethUtils.GCModeFlag.GetName(), ethUtils.GCModeFlag.Value, ethUtils.GCModeFlag.Usage)
	flagSet.Set(ethUtils.GCModeFlag.GetName(), ethUtils.GCModeFlag.Value)

	ctx := cli.NewContext(app, flagSet, nil)
	return ctx
}
