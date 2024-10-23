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
		input, err := readZkProofInput(inputFile)
		if err != nil {
			logger.Failf("Failed to read input file: %v", err)
			return
		}
		proof, err := core.NewZk().ZkProof(input)
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
	// zkProofCmd.Flags().StringVarP(&birthDateThreshold, "threshold", "t", "", "birth date threshold, format is yearmonthdate, such as 20240101")
	// zkProofCmd.Flags().StringVarP(&key, "key", "k", "", "key")
	// zkProofCmd.Flags().StringVarP(&iv, "iv", "i", "", "iv")
	// zkProofCmd.Flags().StringVarP(&iv, "iv", "i", "", "iv")
	zkProofCmd.Flags().StringVarP(&inputFile, "input", "i", "", "input file path which contain input values vc, key, iv, birthdate threshold")
	// zkProofCmd.MarkFlagRequired("vc")
	// zkProofCmd.MarkFlagRequired("threshold")
	zkProofCmd.MarkFlagRequired("input")
}

func readZkProofInput(inputFile string) (*core.ZkProofInput, error) {
	// read input file and unmarshal json to core.ZkProofInput
	input := core.ZkProofInput{}
	content, err := os.ReadFile(inputFile)
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to read input file")
	}
	if err := json.Unmarshal(content, &input); err != nil {
		return nil, errors.WithMessage(err, "Failed to unmarshal input file")
	}
	return &input, nil
}
