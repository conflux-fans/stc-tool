/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package writer

import (
	"github.com/spf13/cobra"
)

// ownerCmd represents the owner command
var ownerCmd = &cobra.Command{
	Use:   "owner",
	Short: "Owner operations",
	Long:  `Owner operations`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var (
	name string
)

func InitCmds(rootCmd *cobra.Command) {
	rootCmd.AddCommand(ownerCmd)

}
