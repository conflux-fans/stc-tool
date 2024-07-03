/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package file

import (
	"github.com/conflux-fans/storage-cli/encrypt"
	"github.com/conflux-fans/storage-cli/logger"
	"github.com/spf13/cobra"
)

// decryptCmd represents the decrypt command
var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt file",
	Long:  `Decrypt file`,
	Run: func(cmd *cobra.Command, args []string) {
		encryptor, err := encrypt.GetEncryptor(cipher)
		if err != nil {
			logger.Failf("Failed to get encryptor %v", err)
			return
		}

		if outputDirPath == "" {
			outputDirPath = "./"
		}

		outputFile, err := encrypt.DecryptFile(encryptor, sourceFilePath, outputDirPath, []byte(password))
		if err != nil {
			logger.Failf("Failed to encrypt file %v", err)
			return
		}

		logger.SuccessWithResult(outputFile, "Decrypted file to below path")
	},
}

func init() {
	fileCmd.AddCommand(decryptCmd)
}
