/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var (
	tag  string
	name string
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload file or data",
	Long:  `Upload file or data`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
	// uploadCmd.Flags().StringVar(&tag, "tag", "", "file tag, for appending content to file")
}
