/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package template

import (
	"github.com/conflux-fans/storage-cli/core"
	"github.com/conflux-fans/storage-cli/logger"
	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download template",
	Long:  `Download template`,
	Run: func(cmd *cobra.Command, args []string) {
		filePath, err := core.DownloadTemplate(name)
		if err != nil {
			logger.Get().WithError(err).Error("Failed to download template")
			return
		}
		logger.Successf("Template file is saved to %s\n", filePath)
	},
}

func init() {
	templateCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().StringVar(&name, "name", "", "template name")
	downloadCmd.MarkFlagRequired("name")
}
