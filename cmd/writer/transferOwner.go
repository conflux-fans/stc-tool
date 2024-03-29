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

// transferOwnerCmd represents the transferOwner command
var transferOwnerCmd = &cobra.Command{
	Use:   "transfer",
	Short: "Transfer content ownership",
	Long:  `Transfer content ownership`,
	Run: func(cmd *cobra.Command, args []string) {
		if !common.IsHexAddress(currentOwner) {
			fmt.Println("from is not valid address")
			return
		}
		if !common.IsHexAddress(newOwner) {
			fmt.Println("to is not valid address")
			return
		}
		err := core.TransferWriter(name, common.HexToAddress(currentOwner), common.HexToAddress(newOwner))
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
	ownerCmd.AddCommand(transferOwnerCmd)
	transferOwnerCmd.Flags().StringVar(&name, "name", "", "content name to transfer ownership")
	transferOwnerCmd.Flags().StringVar(&currentOwner, "from", "", "current owner")
	transferOwnerCmd.Flags().StringVar(&newOwner, "to", "", "target owner")
	transferOwnerCmd.MarkFlagRequired("name")
	transferOwnerCmd.MarkFlagRequired("from")
	transferOwnerCmd.MarkFlagRequired("to")
}
