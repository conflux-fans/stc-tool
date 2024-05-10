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

// contentOwnerCmd represents the contentOwner command
var contentOwnerCmd = &cobra.Command{
	Use:   "content",
	Short: "Check if the account is content owner",
	Long:  `Check if the account is content owner`,
	Run: func(cmd *cobra.Command, args []string) {
		if !common.IsHexAddress(account) {
			logger.Failf("account %v is not valid address", account)
			return
		}

		isWriter, err := core.CheckIsContentWriter(name, common.HexToAddress(account))
		if err != nil {
			logger.Fail(err.Error())
			return
		}

		logger.SuccessfWithParams(map[string]string{
			"Account":      account,
			"Content Name": name,
			"Result":       fmt.Sprintf("%v", isWriter),
		}, "Check if account is owner of content completed")
	},
}

func init() {
	ownerCmd.AddCommand(contentOwnerCmd)
	contentOwnerCmd.Flags().StringVar(&account, "account", "", "account to check if is content owner")
	contentOwnerCmd.Flags().StringVar(&name, "name", "", "content name")
	contentOwnerCmd.MarkFlagRequired("account")
	contentOwnerCmd.MarkFlagRequired("name")
}
