/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package owner

import (
	"github.com/conflux-fans/storage-cli/core"
	"github.com/conflux-fans/storage-cli/logger"
	"github.com/spf13/cobra"
)

// ownerHistoryCmd represents the ownerHistory command
var ownerHistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "query content owner history",
	Long:  `query content owner history`,
	Run: func(cmd *cobra.Command, args []string) {
		history, err := core.DefaultOwnerOperator().GetOwnerHistory(name)
		if err != nil {
			logger.Fail(err.Error())
			return
		}
		logger.SuccessfWithList(history, "Owner history")
	},
}

func init() {
	ownerCmd.AddCommand(ownerHistoryCmd)
	ownerHistoryCmd.Flags().StringVar(&name, "name", "", "content name to query ownership history")
}
