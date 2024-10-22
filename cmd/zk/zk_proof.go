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
		proof, err := core.NewZk().ZkProof(vc, key, iv, birthDateThreshold)
		if err != nil {
			logger.Failf("Failed to generate zk proof: %v", err)
			return
		}

		logger.SuccessfWithParams(map[string]string{
			"VC Proof":  proof.Proof,
			"Flow Root": proof.FlowRoot.Hex(),
		}, "Successfully generated zk proof")
	},
}

var (
	vc                 string
	birthDateThreshold string
	key                string
	iv                 string
	sourceFile         string
)

func init() {
	zkCmd.AddCommand(zkProofCmd)
	// zkProofCmd.Flags().StringVarP(&vc, "vc", "v", "", "vc string in json format")
	// zkProofCmd.Flags().StringVarP(&birthDateThreshold, "threshold", "t", "", "birth date threshold, format is yearmonthdate, such as 20240101")
	// zkProofCmd.Flags().StringVarP(&key, "key", "k", "", "key")
	// zkProofCmd.Flags().StringVarP(&iv, "iv", "i", "", "iv")
	zkProofCmd.Flags().StringVarP(&iv, "iv", "i", "", "iv")
	zkProofCmd.Flags().StringVarP(&sourceFile, "key", "k", "", "key")
	zkProofCmd.MarkFlagRequired("vc")
	zkProofCmd.MarkFlagRequired("threshold")
}
