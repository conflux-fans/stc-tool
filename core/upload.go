package core

import (
	"fmt"

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

// Upload file and create stream with KEY0 by file root hash
func UploadDataByKv(name string, data string) error {
	// get strem by name hash
	// streamId := streamIdByName(name)
	logrus.WithField("name", name).Info("Ready to upload data")
	// iter := kvClientForIterator.NewIterator(STREAM_FILE)
	// revert if exists
	// if err := iter.SeekToFirst(); err != nil {
	// 	return errors.WithMessage(err, "Failed to seek to first")
	// }
	// if iter.Valid() {
	// 	return errors.New("The name already exists")
	// }

	// revert if exists
	if err := checkDataNameExists(name); err != nil {
		return err
	}

	// otherwise upload file
	// if err := UploadFile(filepath, opt); err != nil {
	// 	return errors.WithMessage(err, "Failed to upload file")
	// }
	// write stream with key0 be file root hash

	// rootHash, err := GetRootHash(filepath)
	// if err != nil {
	// 	return errors.WithMessage(err, "Failed to get root hash")
	// }

	// logrus.WithField("root", rootHash).Info("Upload file completed")

	// batcher := kvClientForPut.Batcher()
	// batcher.Set(STREAM_FILE,
	// 	[]byte(getStreamKey(0)),
	// 	[]byte(rootHash[:]),
	// )

	// err = batcher.Exec()
	// if err != nil {
	// 	return errors.WithMessage(err, "Failed to execute batcher")
	// }

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
	v, err := kvClientForIterator.GetValue(STREAM_FILE, []byte(fmt.Sprintf("%s:line", name)))
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

// func getStreamLastKeyId(streamName string) (uint, error) {
// 	streamId := crypto.Keccak256Hash([]byte(streamName))
// 	iter := kvClientForPut.NewIterator(streamId)
// 	if err := iter.SeekToLast(); err != nil {
// 		return 0, err
// 	}
// 	key := string(iter.KeyValue().Key)
// 	return parseStreamKey(key)
// }

// func getStreamKey(id uint) string {
// 	return fmt.Sprintf("FILE%d", id)
// }

// func parseStreamKey(key string) (uint, error) {
// 	idStr := strings.Replace(key, "FILE", "", 1)
// 	idInt, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return uint(idInt), nil
// }

// // Note: useless
// func SaveFileKeyToDb(filepath string, name string) error {
// 	// save db
// 	fileNameKey := db.KeyFileName(name)
// 	_, err := db.GetDB().Get([]byte(fileNameKey), nil)
// 	if err == nil {
// 		return errors.New("already existed")
// 	}
// 	if err != leveldb.ErrNotFound {
// 		return errors.WithMessagef(err, "Failed to query %s", name)
// 	}

// 	rootHash, err := GetRootHash(filepath)
// 	if err != nil {
// 		return err
// 	}
// 	j, err := json.Marshal([]common.Hash{rootHash})
// 	if err != nil {
// 		return err
// 	}

// 	err = db.GetDB().Put([]byte(name), j, nil)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
