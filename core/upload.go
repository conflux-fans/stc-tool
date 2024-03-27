package core

import (
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	ccore "github.com/zero-gravity-labs/zerog-storage-client/core"
	"github.com/zero-gravity-labs/zerog-storage-client/transfer"
	"github.com/zero-gravity-labs/zerog-storage-tool/utils/encryptutils"
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
func UploadDataByKv(name string, data string) error {
	logrus.WithField("name", name).Info("Ready to upload data")

	// revert if exists
	if err := checkDataNameExists(name); err != nil {
		return err
	}

	if err := AppendData(name, data, true); err != nil {
		return err
	}

	logrus.WithField("name", name).Info("Upload data completed")
	return nil
}

func UploadDataFromFile(name string, filePath string) error {
	if err := checkDataNameExists(name); err != nil {
		return err
	}

	if err := AppendFromFile(name, filePath, true); err != nil {
		return err
	}

	logrus.WithField("name", name).Info("Upload data completed")
	return nil
}

func checkDataNameExists(name string) error {
	// revert if exists
	v, err := kvClientForIterator.GetValue(STREAM_FILE, []byte(keyLineCount(name)))
	if err != nil {
		return errors.WithMessage(err, "Failed to get file line size")
	}
	if v.Size != 0 {
		return errors.New("The name already exists")
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
