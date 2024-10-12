/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package zk

import (
	"fmt"

	"github.com/conflux-fans/storage-cli/core"
	"github.com/conflux-fans/storage-cli/logger"
	"github.com/spf13/cobra"
)

// zkVerifyCmd represents the zkVerify command
var zkVerifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "zk verify",
	Long:  `zk verify`,
	Run: func(cmd *cobra.Command, args []string) {
		result, err := core.NewZk().ZkVerify(proof, birthDateThreshold, root)
		if err != nil {
			logger.Failf("Failed verify: %v", err.Error())
			return
		}

		logger.SuccessWithResult(fmt.Sprintf("%v", result), "Verified successfully")
	},
}

var (
	proof, root string
)

func init() {
	zkCmd.AddCommand(zkVerifyCmd)
	zkVerifyCmd.Flags().StringVarP(&proof, "proof", "", "", "proof")
	zkVerifyCmd.Flags().StringVarP(&birthDateThreshold, "birth_threshold", "", "", "birth date threshold")
	zkVerifyCmd.Flags().StringVarP(&root, "root", "", "", "root hash")
}
