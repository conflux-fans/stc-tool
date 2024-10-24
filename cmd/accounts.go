/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/conflux-fans/storage-cli/core"
	"github.com/conflux-fans/storage-cli/logger"
	"github.com/spf13/cobra"
)

// accountsCmd represents the accounts command
var accountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "List all account addresses",
	Long:  `List all account addresses`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.SuccessfWithList(core.GetAccounts(), "all accounts")
	},
}

func init() {
	rootCmd.AddCommand(accountsCmd)
}
