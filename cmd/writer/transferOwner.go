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

// transferOwnerCmd represents the transferOwner command
var transferOwnerCmd = &cobra.Command{
	Use:   "transfer",
	Short: "Transfer content ownership",
	Long:  `Transfer content ownership`,
	Run: func(cmd *cobra.Command, args []string) {
		if !common.IsHexAddress(currentOwner) {
			logger.Failf("from %v is not valid address", currentOwner)
			return
		}
		if !common.IsHexAddress(newOwner) {
			logger.Failf("to %v is not valid address", newOwner)
			return
		}

		core.GrantAllAccountStreamWriter()

		err := core.TransferWriter(name, common.HexToAddress(currentOwner), common.HexToAddress(newOwner))
		if err != nil {
			logger.Fail(err.Error())
			return
		}
		logger.SuccessfWithParams(map[string]string{
			"Name": name,
			"From": currentOwner,
			"To":   newOwner,
		}, "Successfully transferred content ownership")
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
