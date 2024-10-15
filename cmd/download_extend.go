/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/conflux-fans/storage-cli/core"
	"github.com/conflux-fans/storage-cli/logger"
	"github.com/spf13/cobra"
)

// downloadExtendCmd represents the downloadData command
var downloadExtendCmd = &cobra.Command{
	Use:   "content",
	Short: "Download content",
	Long:  `Download content`,
	Run: func(cmd *cobra.Command, args []string) {
		isOutputToconsole, _ := cmd.Flags().GetBool("console")
		isShowMetadata, _ := cmd.Flags().GetBool("metadata")
		if err := core.DefaultDownloader().DownloadExtend(name, isShowMetadata, isOutputToconsole); err != nil {
			logger.Failf("Failed to download content %s: %v\n", name, err)
			return
		}
		logger.SuccessWithResult(name, "Download content successfully")
	},
}

func init() {
	downloadCmd.AddCommand(downloadExtendCmd)
	downloadExtendCmd.Flags().StringVarP(&name, "name", "n", "", "data name")
	downloadExtendCmd.Flags().Bool("console", false, "Output to console")
	downloadExtendCmd.Flags().Bool("metadata", false, "Output metadata")
	downloadExtendCmd.MarkFlagRequired("name")
}
