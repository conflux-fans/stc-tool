/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/zero-gravity-labs/zerog-storage-tool/core"
)

// uploadFileCmd represents the uploadFile command
var uploadFileCmd = &cobra.Command{
	Use:   "file",
	Short: "Upload as file",
	Long:  `Upload as file`,
	Run: func(cmd *cobra.Command, args []string) {
		opt, err := core.NewUploadOption(cipher, password)

		if err != nil {
			logrus.WithError(err).Error("Failed to create encryption option")
			return
		}
		if err := core.UploadFile(filePath, opt); err != nil {
			logrus.WithError(err).Error("Failed to upload file")
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
