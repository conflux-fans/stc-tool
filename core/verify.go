package core

import (
	"os"

	ccore "github.com/0glabs/0g-storage-client/core"
	"github.com/conflux-fans/storage-cli/logger"
	"github.com/conflux-fans/storage-cli/utils/encryptutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

func Verify(filePath string, opt *EncryptOption) (bool, error) {
	// calc file root hash
	if opt != nil {
		outPath, err := encryptutils.EncryptFile(filePath, opt.Method, opt.Password)
		if err != nil {
			return false, errors.WithMessage(err, "Failed to encrypt file")
		}
		filePath = outPath
		defer func() {
			os.Remove(outPath)
		}()
	}

	rootHash, err := GetRootHash(filePath)
	if err != nil {
		return false, err
	}

	logger.Get().WithField("root", rootHash).Info("Data merkle root calculated")

	// get file info by root
	fi, err := GetFileInfo(rootHash)
	if err != nil {
		return false, errors.WithMessage(err, "Failed to get file info from remote")
	}
	if fi == nil {
		logger.Get().Info("Document verification failed due to file not found on node")
		return false, nil
	}

	if fi.Finalized {
		logger.Get().Info("Document verification passed!")
		return true, nil
	}

	logger.Get().Info("Document verification failed due to file upload is not finalized")
	return false, nil
}

func GetRootHash(filePath string) (common.Hash, error) {
	f, err := ccore.Open(filePath)
	if err != nil {
		return common.Hash{}, errors.WithMessage(err, "Failed to open file")
	}

	tree, err := ccore.MerkleTree(f)
	if err != nil {
		return common.Hash{}, errors.WithMessage(err, "Failed to calculate merkel tree root hash")
	}
	return tree.Root(), nil
}
