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
	Short: "batch upload texts",
	Long:  `batch upload texts`,
	Run: func(cmd *cobra.Command, args []string) {
		hash, err := core.BatchUpload(count, &core.EncryptOption{cipher, password})
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// batchuploadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// batchuploadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
