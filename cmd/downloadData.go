/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/conflux-fans/storage-cli/core"
	"github.com/conflux-fans/storage-cli/logger"
	"github.com/spf13/cobra"
)

// downloadDataCmd represents the downloadData command
var downloadDataCmd = &cobra.Command{
	Use:   "data",
	Short: "Download content",
	Long:  `Download content`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := core.DownloadDataByKv(name); err != nil {
			logger.Failf("Failed to download content %s: %v\n", name, err)
		}
	},
}

func init() {
	downloadCmd.AddCommand(downloadDataCmd)
	downloadDataCmd.Flags().StringVarP(&name, "name", "n", "", "data name")
}
