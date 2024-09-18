/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	ccore "github.com/0glabs/0g-storage-client/core"
	"github.com/conflux-fans/storage-cli/constants/enums"
	"github.com/conflux-fans/storage-cli/core"
	"github.com/conflux-fans/storage-cli/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

// uploadExtendCmd represents the uploadData command
var uploadExtendCmd = &cobra.Command{
	Use:   "content",
	Short: "Upload content",
	Long:  `Upload content`,
	Run: func(cmd *cobra.Command, args []string) {
		if !common.IsHexAddress(account) {
			logger.Failf("account %v is not valid address", account)
			return
		}

		var dataType enums.ExtendDataType
		var data ccore.IterableData
		var err error
		if content != "" {
			dataType, data, err = core.DefaultExtendDataConverter().ByContent([]byte(content))
			if err != nil {
				logger.Fail(err.Error())
				return
			}
		}

		if fileOfContent != "" {
			dataType, data, err = core.DefaultExtendDataConverter().ByFile(fileOfContent)
			if err != nil {
				logger.Fail(err.Error())
				return
			}
		}

		if err := core.DefaultUploader().UploadExtend(common.HexToAddress(account), name, dataType, data); err != nil {
			logger.Fail(err.Error())
			return
		}

		logger.SuccessfWithParams(map[string]string{
			"Name": name,
		}, "Upload content completed")
	},
}

func init() {
	uploadCmd.AddCommand(uploadExtendCmd)
	uploadExtendCmd.Flags().StringVar(&fileOfContent, "file", "", "file path of content to upload")
	uploadExtendCmd.Flags().StringVar(&content, "content", "", "content be uploaded")
	uploadExtendCmd.Flags().StringVar(&name, "name", "", "name, for appending content")
	uploadExtendCmd.Flags().StringVar(&account, "account", "", "name, for appending content")
	uploadExtendCmd.MarkFlagsOneRequired("content", "file")
}
