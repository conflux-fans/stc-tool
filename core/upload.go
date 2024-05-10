package core

import (
	"fmt"
	"os"

	ccore "github.com/0glabs/0g-storage-client/core"
	"github.com/0glabs/0g-storage-client/transfer"
	"github.com/conflux-fans/storage-cli/logger"
	"github.com/conflux-fans/storage-cli/utils/encryptutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

type UploadOption struct {
	EncryptOption *EncryptOption
	// Tag           string
}

func NewUploadOption(method string, password string) (*UploadOption, error) {
	encryptOpt, err := NewEncryptOption(method, password)
	if err != nil {
		return nil, err
	}
	return &UploadOption{
		EncryptOption: encryptOpt,
		// Tag:           tag,
	}, nil
}

// Upload data
func UploadDataByKv(account common.Address, name string, data string) error {
	logger.Get().WithField("name", name).Info("Ready to upload data")

	// revert if exists
	if err := checkDataNameExists(name); err != nil {
		return err
	}

	if err := appendData(account, name, data, true); err != nil {
		return err
	}

	logger.Get().WithField("name", name).Info("Upload data completed")
	return nil
}

func UploadDataFromFile(account common.Address, name string, filePath string) error {
	if err := checkDataNameExists(name); err != nil {
		return err
	}

	if err := appendFromFile(account, name, filePath, true); err != nil {
		return err
	}

	logger.Get().WithField("name", name).Info("Upload data completed")
	return nil
}

func checkDataNameExists(name string) error {
	// revert if exists
	v, err := kvClientForIterator.GetValue(STREAM_FILE, []byte(keyLineCount(name)))
	if err != nil {
		return errors.WithMessage(err, "Failed to get file line size")
	}
	if v.Size != 0 {
		return fmt.Errorf("the name %s already exists", name)
	}
	return nil
}

func UploadFile(filepath string, opt *UploadOption) error {
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

	uploader := transfer.NewUploader(defaultFlow, nodeClients)

	f, err := ccore.Open(filepath)
	if err != nil {
		return err
	}

	err = uploader.Upload(f, transfer.UploadOption{
		// Tags: []byte(opt.Tag),
	})

	return err
}
