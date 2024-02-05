/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify file",
	Long:  `Verify file`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var file string

func init() {
	rootCmd.AddCommand(verifyCmd)
	verifyCmd.Flags().StringVar(&cipher, "cipher", "", "cipher method")
	verifyCmd.Flags().StringVar(&password, "password", "", "cipher password")
	verifyCmd.Flags().StringVar(&file, "file", "", "file path")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// verifyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// verifyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
