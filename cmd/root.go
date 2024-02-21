/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/zero-gravity-labs/zerog-storage-tool/cmd/file"
	"github.com/zero-gravity-labs/zerog-storage-tool/cmd/template"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "zerog-storage-tool",
	Short: "zero storage tool",
	Long:  `zerog storage tool for upload,batchupload,append content,download,verify,transfer owner,template manager`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	file.InitCmds(rootCmd)
	template.InitCmds(rootCmd)
}
