/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package version

import (
	"fmt"
	"runtime"

	"fabricbypeer/common/metadata"

	"github.com/spf13/cobra"
)

// Program name
const ProgramName = "peer"

// Cmd returns the Cobra Command for Version
// 返回  子命令 对象
func Cmd() *cobra.Command {
	return cobraCommand
}

//  声明 子对象
var cobraCommand = &cobra.Command{
	Use:   "version",
	Short: "Print fabric peer version.",
	Long:  `Print current version of the fabric peer server.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 判断入参
		if len(args) != 0 {
			return fmt.Errorf("trailing args detected")
		}
		//
		// Parsing of the command line is done so silence cmd usage
		cmd.SilenceUsage = true
		// 输出 GetInfo 函数
		fmt.Print(GetInfo())

		fmt.Println("-------@调试@-40000010001")
		return nil
	},
}

// GetInfo returns version information for the peer
// 暑促胡版本信息
func GetInfo() string {
	ccinfo := fmt.Sprintf("  Base Docker Namespace: %s\n"+
		"  Base Docker Label: %s\n"+
		"  Docker Namespace: %s\n",
		metadata.BaseDockerNamespace,
		metadata.BaseDockerLabel,
		metadata.DockerNamespace)

	return fmt.Sprintf("%s:\n Version: %s\n Commit SHA: %s\n Go version: %s\n"+
		" OS/Arch: %s\n"+
		" Chaincode:\n%s\n",
		ProgramName, metadata.Version, metadata.CommitSHA, runtime.Version(),
		fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH), ccinfo)
}
