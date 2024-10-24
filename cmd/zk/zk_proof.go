/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package zk

import (
	"encoding/json"
	"os"

	"github.com/conflux-fans/storage-cli/core"
	"github.com/conflux-fans/storage-cli/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// zkProofCmd represents the zkProof command
var zkProofCmd = &cobra.Command{
	Use:   "proof",
	Short: "generate zk proof",
	Long:  `generate zk proof`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Get().Info("ready to generate zk proof")
		zkUploadOutput, err := readZkUploadOutput(inputFile)
		if err != nil {
			logger.Failf("Failed to read input file: %v", err)
			return
		}
		logger.Get().Info("read input file")

		zkProofInput, err := core.NewZk().GetZkProofInput(zkUploadOutput.Vc, birthDateThreshold, zkUploadOutput)
		if err != nil {
			logger.Failf("Failed to get zk proof input: %v", err)
			return
		}
		logger.Get().Info("generated zk proof input")

		proof, err := core.NewZk().ZkProof(zkProofInput)
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
	inputFile string
)

func init() {
	zkCmd.AddCommand(zkProofCmd)
	// zkProofCmd.Flags().StringVarP(&vc, "vc", "v", "", "vc string in json format")
	zkProofCmd.Flags().StringVarP(&birthDateThreshold, "threshold", "t", "", "birth date threshold, format is yearmonthdate, such as 20000101")
	// zkProofCmd.Flags().StringVarP(&key, "key", "k", "", "key")
	// zkProofCmd.Flags().StringVarP(&iv, "iv", "i", "", "iv")
	// zkProofCmd.Flags().StringVarP(&iv, "iv", "i", "", "iv")
	zkProofCmd.Flags().StringVarP(&inputFile, "input", "i", "", "input file path which contain input values key, iv, submission_tx_hash. the file is auto generated when run upload command")
	// zkProofCmd.MarkFlagRequired("vc")
	// zkProofCmd.MarkFlagRequired("threshold")
	zkProofCmd.MarkFlagRequired("input")
	zkProofCmd.MarkFlagRequired("threshold")
}

func readZkUploadOutput(inputFile string) (*core.ZkUploadOutput, error) {
	// read input file and unmarshal json to core.ZkProofInput
	input := core.ZkUploadOutput{}
	content, err := os.ReadFile(inputFile)
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to read input file")
	}
	if err := json.Unmarshal(content, &input); err != nil {
		return nil, errors.WithMessage(err, "Failed to unmarshal input file")
	}
	return &input, nil
}
