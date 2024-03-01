package core

import (
	"encoding/json"
	"fmt"
	"strconv"

	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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

// Upload file and create stream with KEY0 by file root hash
func UploadByKv(streamName string, filepath string, opt *UploadOption) error {
	// get strem by name hash
	streamId := streamIdByName(streamName)
	logrus.WithField("stream id", streamId).WithField("stream name", streamName).Info("Ready to upload stream")
	iter := kvClientForIterator.NewIterator(streamId)
	// revert if exists
	if err := iter.SeekToFirst(); err != nil {
		return errors.WithMessage(err, "Failed to seek to first")
	}

	if iter.Valid() {
		return errors.New("The name already exists")
	}
	// otherwise upload file
	if err := uploadFile(filepath, opt); err != nil {
		return errors.WithMessage(err, "Failed to upload file")
	}
	// write stream with key0 be file root hash

	rootHash, err := GetRootHash(filepath)
	if err != nil {
		return errors.WithMessage(err, "Failed to get root hash")
	}

	logrus.WithField("root", rootHash).Info("Upload file completed")

	batcher := kvClientForPut.Batcher()
	batcher.Set(streamId,
		[]byte(getStreamKey(0)),
		[]byte(rootHash[:]),
	)

	err = batcher.Exec()
	if err != nil {
		return errors.WithMessage(err, "Failed to execute batcher")
	}

	logrus.WithField("stream name", streamName).WithField("stream id", streamId).Info("Upload stream completed")
	return nil
}

func uploadFile(filepath string, opt *UploadOption) error {
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

func getStreamLastKeyId(streamName string) (uint, error) {
	streamId := crypto.Keccak256Hash([]byte(streamName))
	iter := kvClientForPut.NewIterator(streamId)
	if err := iter.SeekToLast(); err != nil {
		return 0, err
	}
	key := string(iter.KeyValue().Key)
	return parseStreamKey(key)
}

func getStreamKey(id uint) string {
	return fmt.Sprintf("FILE%d", id)
}

func parseStreamKey(key string) (uint, error) {
	idStr := strings.Replace(key, "FILE", "", 1)
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}
	return uint(idInt), nil
}

// Note: useless
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
