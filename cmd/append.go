/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/conflux-fans/storage-cli/constants/enums"
	"github.com/conflux-fans/storage-cli/core"
	"github.com/conflux-fans/storage-cli/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

// appendCmd represents the append command
var appendCmd = &cobra.Command{
	Use:   "append",
	Short: "Append content to specified file",
	Long:  `Append content to specified file`,
	Run: func(cmd *cobra.Command, args []string) {
		if !common.IsHexAddress(account) {
			logger.Failf("account %s is not valid address", account)
			return
		}

		dataType, data, err := getExtendData()
		if err != nil {
			return
		}
		if dataType == enums.EXTEND_DATA_POINTER {
			logger.Fail("Not support pointer data type")
			return
		}

		if err := core.DefaultAppender().AppendExtend(common.HexToAddress(account), name, data); err != nil {
			logger.Fail(err.Error())
			return
		}

		logger.SuccessfWithParams(map[string]string{
			"Name": name,
		}, "Append content completed")
	},
}

func init() {
	rootCmd.AddCommand(appendCmd)
	appendCmd.Flags().StringVar(&filePath, "file", "", "file path of content to upload")
	appendCmd.Flags().StringVar(&content, "data", "", "append content")
	appendCmd.Flags().StringVar(&name, "name", "", "name, for appending content")
	appendCmd.Flags().StringVar(&account, "account", "", "name, for appending content")
	appendCmd.MarkFlagsOneRequired("data", "file")
}
