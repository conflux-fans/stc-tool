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
	count    int
	cipher   string
	password string
)

// batchuploadCmd represents the batchupload command
var batchuploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Batch upload texts",
	Long:  `Batch upload texts`,
	Run: func(cmd *cobra.Command, args []string) {
		opt, err := core.NewEncryptOption(cipher, password)
		if err != nil {
			logrus.WithError(err).Error("Failed to create encryption option")
			return
		}

		hash, err := core.BatchUpload(count, opt)
		if err != nil {
			logrus.Error("Failed uploading:", err)
		} else {
			logrus.WithField("hash", hash).Info("Upload success")
		}
	},
}

func init() {
	batchCmd.AddCommand(batchuploadCmd)
	batchuploadCmd.Flags().IntVarP(&count, "count", "c", 1, "upload count")
	batchuploadCmd.Flags().StringVar(&cipher, "cipher", "", "cipher method")
	batchuploadCmd.Flags().StringVar(&password, "password", "", "cipher password")
}
