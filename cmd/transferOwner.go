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
		if !common.IsHexAddress(currentOwner) {
			fmt.Println("owner is not valid address")
			return
		}
		err := core.TransferOwner(name, common.HexToAddress(currentOwner), common.HexToAddress(newOwner))
		if err != nil {
			fmt.Println("Failed to transfer owner: ", err)
			return
		}
		fmt.Printf("Successfully transferred ownership of content %s from %v to %v\n", name, currentOwner, newOwner)
	},
}
var (
	currentOwner string
	newOwner     string
)

func init() {
	rootCmd.AddCommand(transferOwnerCmd)
	transferOwnerCmd.PersistentFlags().StringVar(&name, "name", "", "content name to transfer ownership")
	transferOwnerCmd.PersistentFlags().StringVar(&currentOwner, "from", "", "current owner")
	transferOwnerCmd.PersistentFlags().StringVar(&newOwner, "to", "", "new owner")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// transferOwnerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// transferOwnerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
