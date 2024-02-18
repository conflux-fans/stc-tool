package core

import (
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	ccore "github.com/zero-gravity-labs/zerog-storage-client/core"
	"github.com/zero-gravity-labs/zerog-storage-tool/utils/encryptutils"
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

	f, err := ccore.Open(filePath)
	if err != nil {
		return false, errors.WithMessage(err, "Failed to open file")
	}

	tree, err := ccore.MerkleTree(f)
	if err != nil {
		return false, errors.WithMessage(err, "Failed to calculate merkel tree root hash")
	}
	logrus.WithField("root", tree.Root()).Info("Data merkle root calculated")

	// get file info by root
	fi, err := GetFileInfo(tree.Root())
	if err != nil {
		return false, errors.WithMessage(err, "Failed to get file info from remote")
	}
	if fi == nil {
		logrus.Info("Document verification failed due to file not found on node")
		return false, nil
	}

	if fi.Finalized {
		logrus.Info("Document verification passed!")
		return true, nil
	}

	logrus.Info("Document verification failed due to file upload is not finalized")
	return false, nil
}
