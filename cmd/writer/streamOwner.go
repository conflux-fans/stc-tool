/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package writer

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"github.com/zero-gravity-labs/zerog-storage-tool/core"
)

// streamOwnerCmd represents the streamOwner command
var streamOwnerCmd = &cobra.Command{
	Use:   "stream",
	Short: "Check if the account is stream owner",
	Long:  `Check if the account is stream owner`,
	Run: func(cmd *cobra.Command, args []string) {
		if !common.IsHexAddress(account) {
			fmt.Println("account is not valid address")
			return
		}
		isWriter := core.CheckIsStreamWriter(common.HexToAddress(account))
		fmt.Printf("Account %v is stream writer ? %v\n", account, isWriter)
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
