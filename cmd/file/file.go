/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package file

import (
	"github.com/spf13/cobra"
)

// fileCmd represents the query command
var fileCmd = &cobra.Command{
	Use:   "file",
	Short: "File operations",
	Long:  `File operations`,
}

var (
	cipher         string
	password       string
	sourceFilePath string
	outputDirPath  string
)

func InitFileCmds(rootCmd *cobra.Command) {
	rootCmd.AddCommand(fileCmd)
	{
		fileCmd.AddCommand(encryptCmd)
		fileCmd.AddCommand(decryptCmd)
		fileCmd.AddCommand(infoCmd)
	}

	fileCmd.PersistentFlags().StringVar(&cipher, "cipher", "", "cipher method")
	fileCmd.PersistentFlags().StringVar(&password, "password", "", "cipher password")
	fileCmd.PersistentFlags().StringVar(&sourceFilePath, "source", "", "source file path")
	fileCmd.PersistentFlags().StringVar(&outputDirPath, "output", "", "output directory path")
}
