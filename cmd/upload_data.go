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

// uploadByNameCmd represents the uploadData command
var uploadByNameCmd = &cobra.Command{
	Use:   "content",
	Short: "Upload content",
	Long:  `Upload content`,
	Run: func(cmd *cobra.Command, args []string) {
		if !common.IsHexAddress(account) {
			logger.Failf("account %v is not valid address", account)
			return
		}

		if content != "" {
			dataType, data, err := core.DefaultExtendDataConverter().ByContent([]byte(content))
			if err != nil {
				logger.Fail(err.Error())
				return
			}
			if err := core.DefaultUploader().UploadByName(common.HexToAddress(account), name, dataType, data); err != nil {
				logger.Fail(err.Error())
				return
			}
		}

		if fileOfContent != "" {
			dataType, data, err := core.DefaultExtendDataConverter().ByFile(fileOfContent)
			if err != nil {
				logger.Fail(err.Error())
				return
			}
			if err := core.DefaultUploader().UploadByName(common.HexToAddress(account), name, dataType, data); err != nil {
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
	uploadCmd.AddCommand(uploadByNameCmd)
	uploadByNameCmd.Flags().StringVar(&fileOfContent, "file", "", "file path of content to upload")
	uploadByNameCmd.Flags().StringVar(&content, "content", "", "content be uploaded")
	uploadByNameCmd.Flags().StringVar(&name, "name", "", "name, for appending content")
	uploadByNameCmd.Flags().StringVar(&account, "account", "", "name, for appending content")
	uploadByNameCmd.MarkFlagsOneRequired("content", "file")
}
