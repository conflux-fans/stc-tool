package core

import (
	"encoding/json"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
	ccore "github.com/zero-gravity-labs/zerog-storage-client/core"
	"github.com/zero-gravity-labs/zerog-storage-client/transfer"
	"github.com/zero-gravity-labs/zerog-storage-tool/db"
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

	err = uploader.Upload(f, transfer.UploadOption{
		// Tags: []byte(opt.Tag),
	})

	return err
}

func SaveFileKeyToDb(filepath string, name string) error {
	// save db
	fileNameKey := db.KeyFileName(name)
	_, err := db.GetDB().Get([]byte(fileNameKey), nil)
	if err == nil {
		return errors.New("already existed")
	}
	if err != leveldb.ErrNotFound {
		return errors.WithMessagef(err, "Failed to query %s", name)
	}

	rootHash, err := GetRootHash(filepath)
	if err != nil {
		return err
	}
	j, err := json.Marshal([]common.Hash{rootHash})
	if err != nil {
		return err
	}

	err = db.GetDB().Put([]byte(name), j, nil)
	if err != nil {
		return err
	}
	return nil
}
