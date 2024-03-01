/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zero-gravity-labs/zerog-storage-tool/core"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download file",
	Long:  `Download file`,
	Run: func(cmd *cobra.Command, args []string) {
		core.DownloadByKv(name)
	},
}

// var root string

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().StringVarP(&name, "name", "r", "", "file merkle root")
}
