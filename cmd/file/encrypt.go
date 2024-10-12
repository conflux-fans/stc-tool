/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package file

import (
	"github.com/conflux-fans/storage-cli/constants/enums"
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
		cipherMethod, err := enums.ParseCipherMethod(cipher)
		if err != nil {
			logger.Failf("Failed to parse cipher method %v", err)
			return
		}

		encryptor, err := encrypt.GetEncryptor(cipherMethod)
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
