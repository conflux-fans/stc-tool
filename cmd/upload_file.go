/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/conflux-fans/storage-cli/constants/enums"
	"github.com/conflux-fans/storage-cli/core"
	"github.com/conflux-fans/storage-cli/logger"
	"github.com/spf13/cobra"
)

// uploadFileCmd represents the uploadFile command
var uploadFileCmd = &cobra.Command{
	Use:   "file",
	Short: "Upload file",
	Long:  `Upload file`,
	Run: func(cmd *cobra.Command, args []string) {
		cipherMethod, err := enums.ParseCipherMethod(cipher)
		if err != nil {
			logger.Failf("Failed to parse cipher method %v", err)
			return
		}

		opt, err := core.NewUploadOption(cipherMethod, password)

		if err != nil {
			logger.Failf("Failed to create encryption option %v", err)
			return
		}
		tree, err := core.DefaultUploader().UploadFile(filePath, opt)
		if err != nil {
			logger.Failf("Failed to upload file %v", err)
			return
		}
		logger.SuccessfWithParams(map[string]string{
			"File": filePath,
			"Root": tree.Root().Hex(),
		}, "Upload file completed")
	},
}

func init() {
	uploadCmd.AddCommand(uploadFileCmd)
	uploadFileCmd.Flags().StringVar(&filePath, "file", "", "file path to upload")
	uploadFileCmd.Flags().StringVar(&cipher, "cipher", "", "cipher method")
	uploadFileCmd.Flags().StringVar(&password, "password", "", "cipher password")
	uploadFileCmd.MarkFlagRequired("file")
}
