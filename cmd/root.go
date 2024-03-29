/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/conflux-fans/storage-cli/cmd/file"
	"github.com/conflux-fans/storage-cli/cmd/template"
	"github.com/conflux-fans/storage-cli/cmd/writer"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "storage-cli",
	Short: "Storage cli",
	Long:  `Storage cli for upload, append, verify, batchupload, owner manager and template manager`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	file.InitCmds(rootCmd)
	template.InitCmds(rootCmd)
	writer.InitCmds(rootCmd)
}
