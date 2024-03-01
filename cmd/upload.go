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
	tag  string
	name string
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload file",
	Long:  `Upload file`,
	Run: func(cmd *cobra.Command, args []string) {
		opt, err := core.NewUploadOption(cipher, password)

		if err != nil {
			logrus.WithError(err).Error("Failed to create encryption option")
			return
		}
		if err := core.UploadByKv(name, filePath, opt); err != nil {
			logrus.WithError(err).Error("Failed to upload file")
			return
		}
		// if err := core.SaveFileKeyToDb(filePath, name); err != nil {
		// 	logrus.WithError(err).Error("Failed to save file key")
		// 	return
		// }
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
	uploadCmd.Flags().StringVar(&cipher, "cipher", "", "cipher method")
	uploadCmd.Flags().StringVar(&password, "password", "", "cipher password")
	uploadCmd.PersistentFlags().StringVar(&filePath, "file", "", "file path")
	uploadCmd.PersistentFlags().StringVar(&name, "name", "", "file name, for appending content to file")
	// uploadCmd.Flags().StringVar(&tag, "tag", "", "file tag, for appending content to file")

	uploadCmd.MarkFlagRequired("file")
}
