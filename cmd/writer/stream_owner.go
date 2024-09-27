/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package writer

import (
	"fmt"

	"github.com/conflux-fans/storage-cli/core"
	"github.com/conflux-fans/storage-cli/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

// streamOwnerCmd represents the streamOwner command
var streamOwnerCmd = &cobra.Command{
	Use:   "stream",
	Short: "Check if the account is stream owner",
	Long:  `Check if the account is stream owner`,
	Run: func(cmd *cobra.Command, args []string) {
		if !common.IsHexAddress(account) {
			logger.Failf("account %v is not valid address", account)
			return
		}
		isWriter, err := core.CheckIsStreamWriter(common.HexToAddress(account))
		if err != nil {
			logger.Fail(err.Error())
			return
		}

		logger.SuccessfWithParams(map[string]string{
			"Account": account,
			"Result":  fmt.Sprintf("%v", isWriter),
		}, "Check if account is writer of stream completed")
	},
}

var (
	account string
)

func init() {
	writerCmd.AddCommand(streamOwnerCmd)
	streamOwnerCmd.Flags().StringVar(&account, "account", "", "account to check if is stream owner")
	streamOwnerCmd.MarkFlagRequired("account")
}
