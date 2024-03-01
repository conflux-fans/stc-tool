/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package template

import (
	"github.com/spf13/cobra"
)

// templateCmd represents the template command
var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Template opertaions",
	Long:  `Template opertaions`,
}

var (
	name string
	keys []string
)

func InitCmds(rootCmd *cobra.Command) {
	rootCmd.AddCommand(templateCmd)
}
