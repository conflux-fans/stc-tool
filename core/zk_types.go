package core

import (
	"github.com/conflux-fans/storage-cli/zkclient"
	"github.com/ethereum/go-ethereum/common"
)

type ZkUploadOutput struct {
	Vc               *zkclient.VC `json:"vc"`
	VcDataRoot       common.Hash  `json:"vc_data_root"`
	SubmissionTxHash common.Hash  `json:"submission_tx_hash"`
	Key              string       `json:"key"`
	IV               string       `json:"iv"`
}

type FlowProofForZk struct {
	VcDataRoot common.Hash   `json:"vc_data_root"`
	FlowRoot   common.Hash   `json:"flow_root"`
	Lemma      []common.Hash `json:"lemma"`
	Path       uint64        `json:"path"`
}

type ZkProofInput struct {
	*FlowProofForZk
	Vc                 *zkclient.VC `json:"vc"`
	BirthdateThreshold string       `json:"birthdate_threshold"`
	Key                string       `json:"key"`
	IV                 string       `json:"iv"`
}
