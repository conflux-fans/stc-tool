package core

import (
	"os"

	"github.com/pkg/errors"
	ccore "github.com/zero-gravity-labs/zerog-storage-client/core"
	"github.com/zero-gravity-labs/zerog-storage-client/transfer"
	"github.com/zero-gravity-labs/zerog-storage-tool/utils/encryptutils"
)

type UploadOption struct {
	EncryptOption *EncryptOption
	Tag           string
}

func NewUploadOption(method string, password string, tag string) (*UploadOption, error) {
	encryptOpt, err := NewEncryptOption(method, password)
	if err != nil {
		return nil, err
	}
	return &UploadOption{
		EncryptOption: encryptOpt,
		Tag:           tag,
	}, nil
}

func Upload(filepath string, opt *UploadOption) error {
	if opt.EncryptOption != nil {
		outPath, err := encryptutils.EncryptFile(filepath, opt.EncryptOption.Method, opt.EncryptOption.Password)
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
	return uploader.Upload(f, transfer.UploadOption{
		Tags: []byte(opt.Tag),
	})
}
