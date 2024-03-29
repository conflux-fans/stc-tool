/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"github.com/zero-gravity-labs/zerog-storage-tool/core"
)

// uploadDataCmd represents the uploadData command
var uploadDataCmd = &cobra.Command{
	Use:   "data",
	Short: "Upload as data",
	Long:  `Upload as data`,
	Run: func(cmd *cobra.Command, args []string) {
		if !common.IsHexAddress(account) {
			fmt.Println("account is not valid address")
			return
		}

		if data != "" {
			if err := core.UploadDataByKv(common.HexToAddress(account), name, data); err != nil {
				fmt.Println("Faild to append content:", err)
			}
			return
		}

		if filePath != "" {
			if err := core.UploadDataByKv(common.HexToAddress(account), name, filePath); err != nil {
				fmt.Println("Faild to append content from file:", err)
			}
			return
		}
	},
}

func init() {
	uploadCmd.AddCommand(uploadDataCmd)
	uploadDataCmd.Flags().StringVar(&filePath, "file", "", "file path of content to upload")
	uploadDataCmd.Flags().StringVar(&data, "data", "", "append content")
	uploadDataCmd.Flags().StringVar(&name, "name", "", "name, for appending content")
	uploadDataCmd.Flags().StringVar(&account, "account", "", "name, for appending content")
	uploadDataCmd.MarkFlagsOneRequired("data", "file")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uploadDataCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uploadDataCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
