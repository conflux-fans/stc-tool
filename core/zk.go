package core

import (
	"context"
	"encoding/json"
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

func (z *Zk) ZkProof(vc, key, iv, birthdateThreshold string) (*zkclient.ProveOutput, error) {
	var _vc zkclient.VC
	if err := json.Unmarshal([]byte(vc), &_vc); err != nil {
		return nil, errors.WithMessage(err, "failed to parse vc")
	}

	logrus.WithField("flow length", z.mustGetFlowLength()).Info("ready to upload vc data")

	// key, iv := "verysecretkey123", "uniqueiv12345678"
	vcUploadData := _vc.MustGetUploadText(key, iv)
	submissionTx, dataRoot, err := DefaultUploader().UploadBytes(vcUploadData)
	if err != nil {
		return nil, err
	}
	logger.Get().WithField("flow length", z.mustGetFlowLength()).WithField("submission tx", submissionTx).WithField("data root", dataRoot).Info("VC uploaded successfully")

	flowProof, err := z.getSectorProof(submissionTx)
	if err != nil {
		return nil, err
	}

	if flowProof.Lemma[0] != dataRoot {
		logrus.WithField("flow proof data root", flowProof.Lemma[0]).WithField("data root", dataRoot).Error("flow proof data root not match")
	}

	vcProof, err := z.getVcProof(key, iv, _vc, flowProof, birthdateThreshold)
	if err != nil {
		return nil, err
	}

	return &zkclient.ProveOutput{
		Proof:            vcProof,
		VcUploadTextRoot: dataRoot,
		FlowRoot:         flowProof.Lemma[len(flowProof.Lemma)-1],
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

func (z *Zk) getVcProof(key, iv string, vc zkclient.VC, flowProof *node.FlowProof, birthdateThreshold string) (string, error) {
	pathElementms := flowProof.Lemma[1 : len(flowProof.Lemma)-1]
	pathIndex := z.genVcInputPath(flowProof.Path)
	// pathIndex := zkclient.BoolsToUint64(zkclient.InvertBools(flowProof.Path[1 : len(flowProof.Path)-1]))

	birthdate, err := time.Parse("20060102", birthdateThreshold)
	if err != nil {
		return "", errors.WithMessage(err, "failed to parse birthdate threshold")
	}
	logger.Get().WithField("merkel proof", pathElementms).WithField("path index", pathIndex).Info("Ready to gen vc proof")

	return zkClient.GetProof(zkclient.NewProveInput(key, iv, vc, pathElementms, pathIndex, []zkclient.ExtensionSignal{{Date: &birthdate}}))
}

func (z *Zk) genVcInputPath(flowProofPath []bool) uint64 {
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

func (z *Zk) ZkVerify(vcProof string, birthdateThreshold string, root string) (bool, error) {
	logger.Get().WithField("proof", vcProof).
		WithField("birthdate_threshold", birthdateThreshold).
		WithField("root", root).
		Info("Start zk verify")

	birthdate, err := time.Parse("20060102", birthdateThreshold)
	if err != nil {
		return false, errors.WithMessage(err, "failed to parse birthdate threshold")
	}
	return zkClient.Verify(vcProof, zkclient.VerifyInput{
		Extensions: []zkclient.ExtensionSignal{{Date: &birthdate}},
		Root:       common.HexToHash(root),
	})
}
