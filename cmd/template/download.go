/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package template

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/zero-gravity-labs/zerog-storage-tool/core"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "download template",
	Long:  `download template`,
	Run: func(cmd *cobra.Command, args []string) {
		filePath, err := core.DownloadTemplate(name)
		if err != nil {
			logrus.WithError(err).Error("Failed to download template")
			return
		}
		fmt.Printf("Template file is saved to %s\n", filePath)
	},
}

func init() {
	templateCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().StringVar(&name, "name", "", "template name")
	downloadCmd.MarkFlagRequired("name")
}
