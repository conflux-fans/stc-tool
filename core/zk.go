package core

import (
	"context"
	"time"

	"github.com/0glabs/0g-storage-client/node"
	"github.com/conflux-fans/storage-cli/logger"
	"github.com/conflux-fans/storage-cli/zkclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Zk struct {
	client *zkclient.Client
}

func NewZk() *Zk {
	return &Zk{
		client: zkClient,
	}
}

type ZkUploadInput struct {
	Vc                 *zkclient.VC `json:"vc"`
	BirthdateThreshold string       `json:"birthdate_threshold"`
}

type ZkUploadOutput struct {
	FlowProofForZk
	SubmissionTxHash common.Hash `json:"submission_tx_hash"`
	Key              string      `json:"key"`
	IV               string      `json:"iv"`
}

// UploadVc upload vc data to storage and get flow proof, the returned flowProof is removed data root and flow root.
func (z *Zk) UploadVc(input *ZkUploadInput, key, iv string) (*ZkUploadOutput, error) {
	logrus.WithField("flow length", z.mustGetFlowLength()).Info("ready to upload vc data")

	// key, iv := randutils.String(16), randutils.String(16) //"verysecretkey123", "uniqueiv12345678"
	vcUploadData := input.Vc.MustGetUploadText(key, iv)
	submissionTx, uploadedDataRoot, err := DefaultUploader().UploadBytes(vcUploadData)
	if err != nil {
		return nil, err
	}
	logger.Get().WithField("flow length", z.mustGetFlowLength()).WithField("submission tx", submissionTx).WithField("data root", uploadedDataRoot).Info("VC uploaded successfully")

	flowProof, err := z.getSectorProof(submissionTx)
	if err != nil {
		return nil, err
	}

	if flowProof.Lemma[0] != uploadedDataRoot {
		logrus.WithField("flow proof data root", flowProof.Lemma[0]).WithField("data root", uploadedDataRoot).Error("flow proof data root not match")
		return nil, errors.New("data root from flow proof not match with data root")
	}

	return &ZkUploadOutput{
		FlowProofForZk:   *convertFlowProofToForZk(flowProof),
		SubmissionTxHash: submissionTx,
		Key:              key,
		IV:               iv,
	}, nil
}

type FlowProofForZk struct {
	VcDataRoot common.Hash   `json:"vc_data_root"`
	FlowRoot   common.Hash   `json:"flow_root"`
	Lemma      []common.Hash `json:"lemma"`
	Path       uint64        `json:"path"`
}

type ZkProofInput struct {
	ZkUploadInput
	FlowProofForZk
	Key string `json:"key"`
	IV  string `json:"iv"`
}

func (z *Zk) ZkProof(input *ZkProofInput) (*zkclient.ProveOutput, error) {
	birthdate, err := time.Parse("20060102", input.BirthdateThreshold)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to parse birthdate threshold")
	}

	vcProof, err := zkClient.GetProof(zkclient.NewProveInput(input.Key, input.IV, *input.Vc, input.Lemma, input.Path, []zkclient.ExtensionSignal{{Date: &birthdate}}))
	if err != nil {
		return nil, err
	}

	return &zkclient.ProveOutput{
		Proof:      vcProof,
		VcDataRoot: input.VcDataRoot,
		FlowRoot:   input.FlowRoot,
	}, nil
}

// get proof
// 1. get sector position by submission tx log
// 2. get sector proof by storage-client
func (z *Zk) getSectorProof(submissionTxHash common.Hash) (*node.FlowProof, error) {
	receipt, err := adminW3Client.Eth.TransactionReceipt(submissionTxHash)
	if err != nil {
		return nil, err
	}

	submit, err := defaultFlow.ParseSubmit(*receipt.Logs[0].ToEthLog())
	if err != nil {
		return nil, err
	}

	time.Sleep(time.Second * 5)

	return zgNodeClients[0].GetSectorProof(context.Background(), submit.StartPos.Uint64(), nil)
}

// func (z *Zk) getVcProof(key, iv string, vc zkclient.VC, flowProofForZk *FlowProofForZk, birthdateThreshold string) (string, error) {
// 	// pathElementms := flowProof.Lemma[1 : len(flowProof.Lemma)-1]
// 	// pathIndex := z.genVcInputPath(flowProof.Path)
// 	// pathIndex := zkclient.BoolsToUint64(zkclient.InvertBools(flowProof.Path[1 : len(flowProof.Path)-1]))

// 	birthdate, err := time.Parse("20060102", birthdateThreshold)
// 	if err != nil {
// 		return "", errors.WithMessage(err, "failed to parse birthdate threshold")
// 	}
// 	logger.Get().WithField("merkel proof", flowProofForZk.Lemma).WithField("path index", flowProofForZk.Path).Info("Ready to gen vc proof")

// 	return zkClient.GetProof(zkclient.NewProveInput(key, iv, vc, flowProofForZk.Lemma, flowProofForZk.Path, []zkclient.ExtensionSignal{{Date: &birthdate}}))
// }

func convertFlowProofToForZk(flowProof *node.FlowProof) *FlowProofForZk {
	fp := &FlowProofForZk{}
	fp.Lemma = flowProof.Lemma[1 : len(flowProof.Lemma)-1]
	fp.Path = genVcInputPath(flowProof.Path)
	fp.VcDataRoot = flowProof.Lemma[0]
	fp.FlowRoot = flowProof.Lemma[len(flowProof.Lemma)-1]
	return fp
}

func genVcInputPath(flowProofPath []bool) uint64 {
	p := zkclient.InvertBools(flowProofPath)
	p = zkclient.ReverseBools(p)
	return zkclient.BoolsToUint64(p)
}

func (z *Zk) mustGetFlowLength() uint64 {
	tree, err := defaultFlow.Tree(nil)
	if err != nil {
		panic(err)
	}
	return tree.CurrentLength.Uint64()
}

func (z *Zk) ZkVerify(vcProof string, birthdateThreshold string, flowRoot string) (bool, error) {
	logger.Get().WithField("proof", vcProof).
		WithField("birthdate_threshold", birthdateThreshold).
		WithField("root", flowRoot).
		Info("Start zk verify")

	birthdate, err := time.Parse("20060102", birthdateThreshold)
	if err != nil {
		return false, errors.WithMessage(err, "failed to parse birthdate threshold")
	}
	return zkClient.Verify(vcProof, zkclient.VerifyInput{
		Extensions: []zkclient.ExtensionSignal{{Date: &birthdate}},
		Root:       common.HexToHash(flowRoot),
	})
}
