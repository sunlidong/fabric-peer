/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fabricbypeer/bccsp/factory"
	"fabricbypeer/internal/peer/chaincode"
	"fabricbypeer/internal/peer/channel"
	"fabricbypeer/internal/peer/common"
	"fabricbypeer/internal/peer/lifecycle"
	"fabricbypeer/internal/peer/node"
	"fabricbypeer/internal/peer/version"
	_ "net/http/pprof"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// The main command describes the service and
// defaults to printing the help message.
var mainCmd = &cobra.Command{Use: "peer"}

func main() {
	// For environment variables.
	// 给环境变量 设置一个前缀
	viper.SetEnvPrefix(common.CmdRoot)
	// 	//会获取所有的环境变量，同时如果设置过了前缀则会自动补全前缀名
	viper.AutomaticEnv()
	// 替换
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	// Define command-line flags that are valid for all peer commands and
	// subcommands.
	mainFlags := mainCmd.PersistentFlags()

	mainFlags.String("logging-level", "", "Legacy logging level flag")
	viper.BindPFlag("logging_level", mainFlags.Lookup("logging-level"))
	mainFlags.MarkHidden("logging-level")

	// 设置环境变量和一些固定参数 以上

	// TODO  获取工厂对象  返回一个全局 bsscp 对象
	cryptoProvider := factory.GetDefault()

	// 添加 version 子命令
	mainCmd.AddCommand(version.Cmd())
	mainCmd.AddCommand(node.Cmd())
	mainCmd.AddCommand(chaincode.Cmd(nil, cryptoProvider))
	mainCmd.AddCommand(channel.Cmd(nil))
	mainCmd.AddCommand(lifecycle.Cmd(cryptoProvider))

	// On failure Cobra prints the usage message and error string, so we only
	// need to exit with a non-0 status
	if mainCmd.Execute() != nil {
		os.Exit(1)
	}
}
