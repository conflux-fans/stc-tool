/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/zero-gravity-labs/zerog-storage-tool/core"
)

var (
	tag string
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "upload a file or text",
	Long:  `upload a file or text`,
	Run: func(cmd *cobra.Command, args []string) {
		opt, err := core.NewUploadOption(cipher, password, tag)

		if err != nil {
			logrus.WithError(err).Error("Failed to create encryption option")
			return
		}
		if err := core.Upload(filePath, opt); err != nil {
			logrus.WithError(err).Error("Failed to upload file")
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
	uploadCmd.Flags().StringVar(&cipher, "cipher", "", "cipher method")
	uploadCmd.Flags().StringVar(&password, "password", "", "cipher password")
	uploadCmd.Flags().StringVar(&filePath, "file", "", "file path")
	uploadCmd.Flags().StringVar(&tag, "tag", "", "file tag, for appending content to file")

	uploadCmd.MarkFlagRequired("file")
}
