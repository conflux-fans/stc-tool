package core

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/zero-gravity-labs/zerog-storage-tool/db"
)

// append source file to dest name
func AppendFile(source string, destName string) error {
	if err := Upload(source, nil); err != nil {
		return err
	}

	return AppendFileKeyToDb(source, destName)
}

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
