/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package file

import (
	"github.com/conflux-fans/storage-cli/encrypt"
	"github.com/conflux-fans/storage-cli/logger"
	"github.com/spf13/cobra"
)

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt file",
	Long:  `Encrypt file`,
	Run: func(cmd *cobra.Command, args []string) {
		encryptor, err := encrypt.GetEncryptor(cipher)
		if err != nil {
			logger.Failf("Failed to get encryptor %v", err)
			return
		}

		_, err = encrypt.EncryptFile(encryptor, sourceFilePath, outputDirPath, []byte(password))
		if err != nil {
			logger.Failf("Failed to encrypt file %v", err)
			return
		}

		logger.Success("Encrypt file completed!")

	},
}

func init() {
	fileCmd.AddCommand(encryptCmd)
}
