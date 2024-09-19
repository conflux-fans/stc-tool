/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"path"

	"github.com/conflux-fans/storage-cli/core"
	"github.com/conflux-fans/storage-cli/logger"
	"github.com/spf13/cobra"
)

// downloadFileCmd represents the downloadFile command
var downloadFileCmd = &cobra.Command{
	Use:   "file",
	Short: "Download file",
	Long:  `Download file`,
	Run: func(cmd *cobra.Command, args []string) {
		savePath := path.Join(".", root+".zg")
		core.DefaultDownloader().DownloadFile(root, savePath)
		logger.SuccessWithResult(savePath, "Download file successfully, please find in below path")
	},
}

var (
	root string
)

func init() {
	downloadCmd.AddCommand(downloadFileCmd)
	downloadFileCmd.Flags().StringVarP(&root, "root", "r", "", "file merkle root")
	downloadFileCmd.MarkFlagRequired("root")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadFileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadFileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
