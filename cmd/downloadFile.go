/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"path"

	"github.com/spf13/cobra"
	"github.com/zero-gravity-labs/zerog-storage-tool/core"
)

// downloadFileCmd represents the downloadFile command
var downloadFileCmd = &cobra.Command{
	Use:   "file",
	Short: "Download file",
	Long:  `Download file`,
	Run: func(cmd *cobra.Command, args []string) {
		savePath := path.Join(".", root+".zg")
		core.DownloadFile(root, savePath)
		fmt.Printf("Download file successfully, please find in %s\n", savePath)
	},
}

var (
	root string
)

func init() {
	downloadCmd.AddCommand(downloadFileCmd)
	downloadFileCmd.Flags().StringVarP(&root, "root", "r", "", "file merkle root")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadFileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadFileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
