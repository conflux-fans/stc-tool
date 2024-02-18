package core

import (
	"os"

	"github.com/pkg/errors"
	ccore "github.com/zero-gravity-labs/zerog-storage-client/core"
	"github.com/zero-gravity-labs/zerog-storage-client/transfer"
	"github.com/zero-gravity-labs/zerog-storage-tool/utils/encryptutils"
)

func Upload(filepath string, opt *EncryptOption) error {
	if opt != nil {
		outPath, err := encryptutils.EncryptFile(filepath, opt.Method, opt.Password)
		if err != nil {
			return errors.WithMessage(err, "Failed to encrypt file")
		}
		filepath = outPath
		defer func() {
			os.Remove(outPath)
		}()
	}

	uploader := transfer.NewUploader(flow, nodeClients)

	f, err := ccore.Open(filepath)
	if err != nil {
		return err
	}
	return uploader.Upload(f)
}
