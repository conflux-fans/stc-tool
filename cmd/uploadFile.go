/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
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
		opt, err := core.NewUploadOption(cipher, password)

		if err != nil {
			logger.Failf("Failed to create encryption option %v", err)
			return
		}
		if err := core.UploadFile(filePath, opt); err != nil {
			logger.Fail(err.Error())
			return
		}
	},
}

func init() {
	uploadCmd.AddCommand(uploadFileCmd)
	uploadFileCmd.Flags().StringVar(&filePath, "file", "", "file path to upload")
	uploadFileCmd.Flags().StringVar(&cipher, "cipher", "", "cipher method")
	uploadFileCmd.Flags().StringVar(&password, "password", "", "cipher password")
	uploadFileCmd.MarkFlagRequired("file")
}
