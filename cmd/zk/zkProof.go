/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package zk

import (
	"github.com/conflux-fans/storage-cli/core"
	"github.com/conflux-fans/storage-cli/logger"
	"github.com/spf13/cobra"
)

// zkProofCmd represents the zkProof command
var zkProofCmd = &cobra.Command{
	Use:   "proof",
	Short: "generate zk proof",
	Long:  `generate zk proof`,
	Run: func(cmd *cobra.Command, args []string) {
		proof, err := core.ZkProof(vc, birthDateThreshold)
		if err != nil {
			logger.Failf("Failed to generate zk proof: %v", err)
			return
		}

		logger.SuccessfWithParams(map[string]string{
			"VC Proof":       proof.Proof,
			"Encrypt VcRoot": proof.EncryptVcRoot.Hex(),
			"Flow Root":      proof.FlowRoot.Hex(),
		}, "Successfully generated zk proof")
	},
}

var (
	vc                 string
	birthDateThreshold string
	// pathElements       string
	// pathIndices        string
)

func init() {
	zkCmd.AddCommand(zkProofCmd)
	zkProofCmd.Flags().StringVarP(&vc, "vc", "v", "", "vc string in json format")
	zkProofCmd.Flags().StringVarP(&birthDateThreshold, "threshold", "t", "", "birth date threshold, format is yearmonthdate, such as 20240101")
	// zkProofCmd.Flags().StringVarP(&pathElements, "path-elements", "e", "", "path elements")
	// zkProofCmd.Flags().StringVarP(&pathIndices, "path-indicies", "i", "", "path indices")
	zkProofCmd.MarkFlagRequired("vc")
	zkProofCmd.MarkFlagRequired("threshold")
}
