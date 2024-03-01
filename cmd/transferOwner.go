/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"github.com/zero-gravity-labs/zerog-storage-tool/core"
)

// transferOwnerCmd represents the transferOwner command
var transferOwnerCmd = &cobra.Command{
	Use:   "transferOwner",
	Short: "Transfer stream admin",
	Long:  `Transfer stream admin`,
	Run: func(cmd *cobra.Command, args []string) {
		if !common.IsHexAddress(owner) {
			fmt.Println("owner is not valid address")
			return
		}
		err := core.TransferOwner(name, common.HexToAddress(owner))
		if err != nil {
			fmt.Println("Failed to transfer owner: ", err)
			return
		}
		fmt.Println("Successfully transferred ownership")
	},
}
var (
	owner string
)

func init() {
	rootCmd.AddCommand(transferOwnerCmd)
	transferOwnerCmd.PersistentFlags().StringVar(&name, "name", "", "file name to transfer ownership")
	transferOwnerCmd.PersistentFlags().StringVar(&owner, "owner", "", "new owner")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// transferOwnerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// transferOwnerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
