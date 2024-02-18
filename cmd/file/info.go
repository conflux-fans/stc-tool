/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package file

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/zero-gravity-labs/zerog-storage-tool/core"
)

// queryCmd represents the query command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get file info by hash",
	Long:  `Get file info by hash`,
	Run: func(cmd *cobra.Command, args []string) {
		fi, err := core.GetFileInfo(common.HexToHash(rootHash))
		if err != nil {
			logrus.WithError(err).Error("Failed to get file info")
		} else {
			logrus.WithField("fi", fi).Info("Get file info completed")
		}
	},
}

var (
	rootHash string
)

func init() {
	infoCmd.Flags().StringVarP(&rootHash, "root", "r", "", "root hash of content")
	infoCmd.MarkFlagRequired("root")
}
