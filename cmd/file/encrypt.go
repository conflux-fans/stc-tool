/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package file

import (
	"fmt"

	"github.com/conflux-fans/storage-cli/encrypt"
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
			fmt.Println("Failed to get encryptor", err)
		}

		_, err = encrypt.EncryptFile(encryptor, sourceFilePath, outputDirPath, []byte(password))
		if err != nil {
			fmt.Println("Failed to encrypt file", err)
		}
	},
}

func init() {
	fileCmd.AddCommand(encryptCmd)
}
