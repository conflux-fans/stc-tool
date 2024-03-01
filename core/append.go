package core

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/zero-gravity-labs/zerog-storage-tool/db"
)

// append source file to dest name
func AppendFile(streamName string, filepath string) error {
	// get strem by name hash
	streamId := streamIdByName(streamName)
	iter := kvClientForPut.NewIterator(streamId)
	// revert if exists
	iter.SeekToLast()
	if !iter.Valid() {
		return errors.New("The name unexists")
	}

	lastKey := string(iter.KeyValue().Key)
	lastFileId, err := parseStreamKey(lastKey)
	if err != nil {
		return err
	}

	// otherwise upload file
	if err := uploadFile(filepath, nil); err != nil {
		return errors.New("Failed to upload file")
	}
	// write stream with key0 be file root hash

	rootHash, err := GetRootHash(filepath)
	if err != nil {
		return errors.WithMessage(err, "Failed to get root hash")
	}

	batcher := kvClientForPut.Batcher()
	batcher.Set(streamId,
		[]byte(getStreamKey(lastFileId+1)),
		[]byte(rootHash[:]),
	)

	err = batcher.Exec()
	if err != nil {
		return errors.WithMessage(err, "Failed to execute batcher")
	}
	return nil
}

// Note: useless
func AppendFileKeyToDb(filepath string, name string) error {
	// save db
	fileNameKey := db.KeyFileName(name)
	value, err := db.GetDB().Get([]byte(fileNameKey), nil)
	if err != nil {
		return errors.WithMessagef(err, "Failed to query %s", name)
	}

	var roots []common.Hash
	if err := json.Unmarshal(value, &roots); err != nil {
		return err
	}

	rootHash, err := GetRootHash(filepath)
	if err != nil {
		return err
	}

	j, err := json.Marshal(append(roots, rootHash))
	if err != nil {
		return err
	}

	err = db.GetDB().Put([]byte(name), j, nil)
	if err != nil {
		return err
	}
	return nil
}
