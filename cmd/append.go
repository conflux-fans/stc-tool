/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// appendCmd represents the append command
var appendCmd = &cobra.Command{
	Use:   "append",
	Short: "append content to specified file",
	Long:  `append content to specified file`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("append called")
	},
}

func init() {
	uploadCmd.AddCommand(appendCmd)
}
