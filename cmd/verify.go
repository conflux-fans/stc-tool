/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/conflux-fans/storage-cli/core"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify file",
	Long:  `Verify file`,
	Run: func(cmd *cobra.Command, args []string) {
		opt, err := core.NewEncryptOption(cipher, password)
		if err != nil {
			logrus.WithError(err).Error("Failed to create options")
			return
		}
		isPassed, err := core.Verify(filePath, opt)
		if err != nil {
			logrus.WithError(err).Error("Failed to verify file")
			return
		}
		logrus.WithField("Passed", isPassed).Info("Document verification completed")

	},
}

var filePath string

func init() {
	rootCmd.AddCommand(verifyCmd)
	verifyCmd.Flags().StringVar(&cipher, "cipher", "", "cipher method")
	verifyCmd.Flags().StringVar(&password, "password", "", "cipher password")
	verifyCmd.Flags().StringVar(&filePath, "file", "", "file path")
}
