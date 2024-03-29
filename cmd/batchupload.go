/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/conflux-fans/storage-cli/core"
	"github.com/spf13/cobra"
)

var (
	count    int
	cipher   string
	password string
)

// batchuploadCmd represents the batchupload command
var batchuploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Batch upload content",
	Long:  `Batch upload content`,
	Run: func(cmd *cobra.Command, args []string) {
		core.BatchUploadByKv(count)
	},
}

func init() {
	batchCmd.AddCommand(batchuploadCmd)
	batchuploadCmd.Flags().IntVarP(&count, "count", "c", 1, "upload count")
	batchuploadCmd.Flags().StringVar(&cipher, "cipher", "", "cipher method")
	batchuploadCmd.Flags().StringVar(&password, "password", "", "cipher password")
}
