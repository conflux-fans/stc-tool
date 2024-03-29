/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package writer

import (
	"fmt"

	"github.com/conflux-fans/storage-cli/core"
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
			fmt.Println("account is not valid address")
			return
		}
		isWriter := core.CheckIsContentWriter(name, common.HexToAddress(account))
		fmt.Printf("Account %v is writer of content %s ? %v\n", account, name, isWriter)
	},
}

func init() {
	ownerCmd.AddCommand(contentOwnerCmd)
	contentOwnerCmd.Flags().StringVar(&account, "account", "", "account to check if is content owner")
	contentOwnerCmd.Flags().StringVar(&name, "name", "", "content name")
	contentOwnerCmd.MarkFlagRequired("account")
	contentOwnerCmd.MarkFlagRequired("name")
}
