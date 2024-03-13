/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package file

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zero-gravity-labs/zerog-storage-tool/encrypt"
)

// decryptCmd represents the decrypt command
var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt file",
	Long:  `Decrypt file`,
	Run: func(cmd *cobra.Command, args []string) {
		encryptor, err := encrypt.GetEncryptor(cipher)
		if err != nil {
			fmt.Println("Failed to get encryptor", err)
			return
		}

		outputFile, err := encrypt.DecryptFile(encryptor, sourceFilePath, outputDirPath, []byte(password))
		if err != nil {
			fmt.Println("Failed to encrypt file", err)
			return
		}

		fmt.Printf("Decrypted file path %s\n", outputFile)
	},
}

func init() {
	fileCmd.AddCommand(decryptCmd)
}
