/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ownerCmd represents the owner command
var ownerCmd = &cobra.Command{
	Use:   "owner",
	Short: "Get content owner",
	Long:  `Get content owner`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("owner called")
	},
}

func init() {
	rootCmd.AddCommand(ownerCmd)
	ownerCmd.Flags().StringVar(&name, "name", "", "content name to check ownership")
	ownerCmd.MarkFlagRequired("name")
}
