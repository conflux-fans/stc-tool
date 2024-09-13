package core

import (
	"encoding/json"

	"github.com/conflux-fans/storage-cli/logger"
	"github.com/conflux-fans/storage-cli/utils/encryptutils"
	"github.com/conflux-fans/storage-cli/zkclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

func ZkProof(vc string, birthdateThreshold string) (*zkclient.ProveOutput, error) {
	var _vc zkclient.VC
	if err := json.Unmarshal([]byte(vc), &_vc); err != nil {
		return nil, errors.WithMessage(err, "failed to parse vc")
	}

	// upload encrypted vc
	encryptedVc, err := encryptutils.EncryptBytes(_vc.EncodeAndPadToSector(), "aes", "1234567812345678")
	if err != nil {
		return nil, errors.WithMessage(err, "failed to encrypt vc")
	}

	segTree, _, err := DefaultUploader().UploadString(zkclient.PadToSector(encryptedVc))
	if err != nil {
		return nil, err
	}
	logger.Get().Info("encrypted vc uploaded successfully")

	// upload vc hash
	data := _vc.Hash()
	_, chunksTree, err := DefaultUploader().UploadString(zkclient.PadToSector(data[:]))
	if err != nil {
		return nil, err
	}
	logger.Get().Info("vc hash uploaded successfully")

	// get vc hash merkel proof against flow
	// segmentWithProof, err := nodeClients[0].ZeroGStorage().DownloadSegmentWithProof(segmentTree.Root(), 0)
	// if err != nil {
	// 	return "", err
	// }

	// lemma := segmentWithProof.Proof.Lemma
	// PathElements := lemma[:len(lemma)-1]

	// pathIndex := zkclient.BoolsToUint64(segmentWithProof.Proof.Path)
	// logger.Get().WithField("path elements", lemma).WithField("path index", pathIndex).WithField("root", segmentTree).Info("Get merkel proof")

	// chunksTree, err := getChunksTree(data)
	// if err != nil {
	// 	return "", err
	// }

	// TODO: gen proof
	lemma := chunksTree.ProofAt(0).Lemma
	PathElements := lemma[:len(lemma)-1]

	pathIndex := zkclient.BoolsToUint64(chunksTree.ProofAt(0).Path)

	vcProof, err := zkClient.GetProof(&zkclient.ProveInput{
		Data:               _vc,
		BirthdateThreshold: birthdateThreshold,
		MerkleProof:        PathElements,
		PathIndex:          pathIndex,
	})
	if err != nil {
		return nil, err
	}

	// get flow root
	flowRoot, err := defaultFlow.Root(nil)
	if err != nil {
		return nil, err
	}

	return &zkclient.ProveOutput{
		Proof:         vcProof,
		EncryptVcRoot: segTree.Root(),
		FlowRoot:      flowRoot,
	}, nil
}

func ZkVerify(vcProof string, birthdateThreshold string, root string) (bool, error) {
	logger.Get().WithField("proof", vcProof).
		WithField("birthdate_threshold", birthdateThreshold).
		// WithField("leaf_hash", leafHash).
		WithField("root", root).
		Info("start zk verify")

	// var _pubInputs []string
	// if err := json.Unmarshal([]byte(pubInputs), &_pubInputs); err != nil {
	// 	return false, errors.WithMessage(err, "failed to unmarshal public inputs")
	// }
	return zkClient.Verify(vcProof, zkclient.VerifyInput{
		BirthdateThreshold: birthdateThreshold,
		Root:               common.HexToHash(root),
		// LeafHash:           common.HexToHash(leafHash),

	})
}
