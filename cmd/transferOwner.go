/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// transferOwnerCmd represents the transferOwner command
var transferOwnerCmd = &cobra.Command{
	Use:   "transferOwner",
	Short: "transfer stream writter role",
	Long:  `transfer stream writter role`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("transferOwner called")
	},
}

func init() {
	rootCmd.AddCommand(transferOwnerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// transferOwnerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// transferOwnerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
