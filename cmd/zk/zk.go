/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package zk

import (
	"fmt"

	"github.com/spf13/cobra"
)

// zkCmd represents the zk command
var zkCmd = &cobra.Command{
	Use:   "zk",
	Short: "generate zk proof and verify",
	Long:  `generate zk proof and verify`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("zk called")
	},
}

func InitCmds(rootCmd *cobra.Command) {
	rootCmd.AddCommand(zkCmd)
}
