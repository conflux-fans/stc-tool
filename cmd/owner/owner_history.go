/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package owner

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ownerHistoryCmd represents the ownerHistory command
var ownerHistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "query content owner history",
	Long:  `query content owner history`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ownerHistory called")
	},
}

func init() {
	ownerCmd.AddCommand(ownerHistoryCmd)
	ownerHistoryCmd.Flags().StringVar(&name, "name", "", "content name to query ownership history")
}
