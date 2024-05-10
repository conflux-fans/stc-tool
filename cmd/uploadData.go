/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/conflux-fans/storage-cli/core"
	"github.com/conflux-fans/storage-cli/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

// uploadDataCmd represents the uploadData command
var uploadDataCmd = &cobra.Command{
	Use:   "content",
	Short: "Upload content",
	Long:  `Upload content`,
	Run: func(cmd *cobra.Command, args []string) {
		if !common.IsHexAddress(account) {
			logger.Failf("account %v is not valid address", account)
			return
		}

		if data != "" {
			if err := core.UploadDataByKv(common.HexToAddress(account), name, data); err != nil {
				logger.Fail(err.Error())
			}
			return
		}

		if filePath != "" {
			if err := core.UploadDataByKv(common.HexToAddress(account), name, filePath); err != nil {
				logger.Fail(err.Error())
			}
			return
		}
	},
}

func init() {
	uploadCmd.AddCommand(uploadDataCmd)
	uploadDataCmd.Flags().StringVar(&filePath, "file", "", "file path of content to upload")
	uploadDataCmd.Flags().StringVar(&data, "content", "", "content be uploaded")
	uploadDataCmd.Flags().StringVar(&name, "name", "", "name, for appending content")
	uploadDataCmd.Flags().StringVar(&account, "account", "", "name, for appending content")
	uploadDataCmd.MarkFlagsOneRequired("content", "file")
}
