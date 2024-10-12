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

// decryptCmd represents the decrypt command
var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt file",
	Long:  `Decrypt file`,
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

		if outputDirPath == "" {
			outputDirPath = "./"
		}

		outputFile, err := encrypt.DecryptFile(encryptor, sourceFilePath, outputDirPath, []byte(password))
		if err != nil {
			logger.Failf("Failed to decrypt file %v", err)
			return
		}

		logger.SuccessWithResult(outputFile, "Decrypted file to below path")
	},
}

func init() {
	fileCmd.AddCommand(decryptCmd)
}
