/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/conflux-fans/storage-cli/constants/enums"
	"github.com/conflux-fans/storage-cli/core"
	"github.com/conflux-fans/storage-cli/logger"
	"github.com/spf13/cobra"
)

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify file",
	Long:  `Verify if file match with providing file`,
	Run: func(cmd *cobra.Command, args []string) {
		cipherMethod, err := enums.ParseCipherMethod(cipher)
		if err != nil {
			logger.Failf("Failed to parse cipher method %v", err)
			return
		}
		opt, err := core.NewEncryptOption(cipherMethod, password)
		if err != nil {
			logger.Failf("Failed to create options: %s", err.Error())
			return
		}
		isPassed, err := core.Verify(filePath, opt)
		if err != nil {
			fmt.Println(err.Error())
			logger.Failf("Failed to verify file: %s", err.Error())
			fmt.Println()
			return
		}
		// logger.Get().WithField("Passed", isPassed).Info("Document verification completed")

		logger.SuccessfWithParams(map[string]string{
			"Is Pass": fmt.Sprintf("%v", isPassed),
		}, "Document verification completed")
	},
}

var filePath string

func init() {
	rootCmd.AddCommand(verifyCmd)
	verifyCmd.Flags().StringVar(&cipher, "cipher", "", "cipher method")
	verifyCmd.Flags().StringVar(&password, "password", "", "cipher password")
	verifyCmd.Flags().StringVar(&filePath, "file", "", "file path")
}
