/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package writer

import (
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
		isWriter := core.CheckIsStreamWriter(common.HexToAddress(account))
		// fmt.Printf("Account %v is stream writer ? %v\n", account, isWriter)
		if isWriter {
			logger.Successf("Account %v is stream writer", account)
		} else {
			logger.Successf("Account %v isn't stream writer", account)
		}
	},
}

var (
	account string
)

func init() {
	ownerCmd.AddCommand(streamOwnerCmd)
	streamOwnerCmd.Flags().StringVar(&account, "account", "", "account to check if is stream owner")
	streamOwnerCmd.MarkFlagRequired("account")
}
