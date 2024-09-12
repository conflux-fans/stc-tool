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

// uploadContentCmd represents the uploadData command
var uploadContentCmd = &cobra.Command{
	Use:   "content",
	Short: "Upload content",
	Long:  `Upload content`,
	Run: func(cmd *cobra.Command, args []string) {
		if !common.IsHexAddress(account) {
			logger.Failf("account %v is not valid address", account)
			return
		}

		if content != "" {
			if err := core.DefaultUploader().UploadDataFromContent(common.HexToAddress(account), name, content); err != nil {
				logger.Fail(err.Error())
				return
			}
		}

		if filePath != "" {
			if err := core.DefaultUploader().UploadDataFromFile(common.HexToAddress(account), name, fileOfContent); err != nil {
				logger.Fail(err.Error())
				return
			}
		}

		logger.SuccessfWithParams(map[string]string{
			"Name": name,
		}, "Upload content completed")
	},
}

func init() {
	uploadCmd.AddCommand(uploadContentCmd)
	uploadContentCmd.Flags().StringVar(&fileOfContent, "file", "", "file path of content to upload")
	uploadContentCmd.Flags().StringVar(&content, "content", "", "content be uploaded")
	uploadContentCmd.Flags().StringVar(&name, "name", "", "name, for appending content")
	uploadContentCmd.Flags().StringVar(&account, "account", "", "name, for appending content")
	uploadContentCmd.MarkFlagsOneRequired("content", "file")
}
