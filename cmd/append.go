/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zero-gravity-labs/zerog-storage-tool/core"
)

// appendCmd represents the append command
var appendCmd = &cobra.Command{
	Use:   "append",
	Short: "Append content to specified file",
	Long:  `Append content to specified file`,
	Run: func(cmd *cobra.Command, args []string) {
		if data != "" {
			if err := core.AppendData(name, data, false); err != nil {
				fmt.Println("Faild to append content:", err)
			}
			return
		}

		if filePath != "" {
			if err := core.AppendFromFile(name, filePath, false); err != nil {
				fmt.Println("Faild to append content from file:", err)
			}
			return
		}
	},
}

var (
	data string
)

func init() {
	rootCmd.AddCommand(appendCmd)
	appendCmd.Flags().StringVar(&filePath, "file", "", "file path of content to upload")
	appendCmd.Flags().StringVar(&data, "data", "", "append content")
	appendCmd.Flags().StringVar(&name, "name", "", "name, for appending content")
	appendCmd.MarkFlagsOneRequired("data", "file")
}
