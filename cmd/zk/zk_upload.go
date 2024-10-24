/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package zk

import (
	"encoding/json"

	"github.com/conflux-fans/storage-cli/core"
	"github.com/conflux-fans/storage-cli/logger"
	"github.com/conflux-fans/storage-cli/utils/randutils"
	"github.com/conflux-fans/storage-cli/zkclient"
	"github.com/spf13/cobra"
)

// zkUploadCmd represents the zkUpload command
var zkUploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "upload vc",
	Long: `upload VC data to the storage system. The output will be automatically saved to a file for the zk proof command, containing the following fields:
- key: The key used for encryption
- iv: The initialization vector
- submission_tx_hash: The hash of the submission transaction
- vc_data_root: The root hash of the VC encrypted data
`,
	Run: func(cmd *cobra.Command, args []string) {
		var _vc zkclient.VC
		err := json.Unmarshal([]byte(vc), &_vc)
		if err != nil {
			logger.Failf("Failed to unmarshal vc: %v", err)
			return
		}

		key, iv := randutils.String(16), randutils.String(16)
		uploadOut, resultPath, err := core.NewZk().UploadVc(&_vc, key, iv)
		if err != nil {
			logger.Failf("Failed to upload vc: %v", err)
			return
		}

		logger.SuccessfWithParams(map[string]string{
			"key":                key,
			"iv":                 iv,
			"vc_data_root":       uploadOut.VcDataRoot.Hex(),
			"submission_tx_hash": uploadOut.SubmissionTxHash.Hex(),
		}, "Successfully uploaded VC. The result is saved to file %v, which can be used as input for the zk proof command.", resultPath)
	},
}

var (
	vc string
)

func init() {
	zkCmd.AddCommand(zkUploadCmd)
	zkUploadCmd.Flags().StringVarP(&vc, "vc", "v", "", "vc string in json format")
	// zkProofCmd.Flags().StringVarP(&birthDateThreshold, "threshold", "t", "", "birth date threshold, format is yearmonthdate, such as 20240101")
	// zkProofCmd.Flags().StringVarP(&key, "key", "k", "", "key")
	// zkProofCmd.Flags().StringVarP(&iv, "iv", "i", "", "iv")

	// zkProofCmd.Flags().StringVarP(&inputFile, "input", "i", "", "input file path which contain input values vc, key, iv, birthdate threshold")
	// zkProofCmd.MarkFlagRequired("vc")
	// zkProofCmd.MarkFlagRequired("threshold")
}
