package core

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/zero-gravity-labs/zerog-storage-client/core"
	"github.com/zero-gravity-labs/zerog-storage-tool/db"
)

var (
	STREAM_FILE = common.HexToHash("000000000000000000000000000000000000000000000000000000000000f2bd")
)

const (
	CHUNK_SIZE     = 4096
	VALUE_MAX_SIZE = CHUNK_SIZE * 100
)

// append source file to dest name
func AppendData(name string, content string, force bool) error {
	if len(content) > VALUE_MAX_SIZE {
		return errors.New("Exceed max size once uploadable")
	}

	// split content to chunks
	var chunks []string
	for i := 0; i < len(content); i += CHUNK_SIZE {
		end := lo.Min([]int{(i + 1) * CHUNK_SIZE, len(content)})
		chunks = append(chunks, content[i*CHUNK_SIZE:end])
	}

	// query size
	v, err := kvClientForIterator.GetValue(STREAM_FILE, []byte(fmt.Sprintf("%s:line", name)))
	if err != nil {
		return errors.WithMessage(err, "Failed to get file line size")
	}

	lineCount := 0
	if v.Size == 0 {
		if !force {
			return errors.New("Unexists name")
		}
	} else {
		// set key to name:line, value to chunk
		lineCountInStr := string(v.Data)
		lineCount, _ = strconv.Atoi(lineCountInStr)
	}

	batcher := kvClientForPut.Batcher()
	batcher.Set(STREAM_FILE,
		[]byte(fmt.Sprintf("%s:line", name)),
		[]byte(fmt.Sprintf("%d", lineCount+len(chunks))),
	)
	for i, chunk := range chunks {
		batcher.Set(STREAM_FILE,
			[]byte(fmt.Sprintf("%s:%d", name, lineCount+i)),
			[]byte(chunk),
		)
	}

	err = batcher.Exec()
	if err != nil {
		return errors.WithMessage(err, "Failed to execute batcher")
	}

	logrus.WithField("line", len(chunks)).Info("Append content completed")

	return nil

	// iter := kvClientForPut.NewIterator(FILE_STREAM)
	// // revert if exists
	// iter.SeekToLast()
	// if !iter.Valid() {
	// 	return errors.New("The name unexists")
	// }

	// lastKey := string(iter.KeyValue().Key)
	// lastFileId, err := parseStreamKey(lastKey)
	// if err != nil {
	// 	return err
	// }

	// otherwise upload file
	// if err := uploadFile(filePath, nil); err != nil {
	// 	return errors.New("Failed to upload file")
	// }
	// // write stream with key0 be file root hash

	// rootHash, err := GetRootHash(filePath)
	// if err != nil {
	// 	return errors.WithMessage(err, "Failed to get root hash")
	// }

	// batcher := kvClientForPut.Batcher()
	// batcher.Set(FILE_STREAM,
	// 	[]byte(getStreamKey(lastFileId+1)),
	// 	[]byte(rootHash[:]),
	// )

	// err = batcher.Exec()
	// if err != nil {
	// 	return errors.WithMessage(err, "Failed to execute batcher")
	// }
	// return nil
}

func AppendFromFile(name string, filePath string, force bool) error {
	// read file
	// if err := UploadFile(filePath, nil); err != nil {
	// 	return errors.New("Failed to upload file")
	// }

	f, err := openFile(filePath)
	if err != nil {
		return err
	}

	// split by VALUE_MAX_SIZE
	for {
		buffer := make([]byte, 1024)
		n, err := f.Read(buffer)
		if err != nil {
			return err
		}
		if n == 0 {
			return nil
		}

		if err = AppendData(name, string(buffer[:n]), force); err != nil {
			return err
		}
	}
}

func openFile(name string) (*os.File, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		return nil, core.ErrFileRequired
	}

	if info.Size() == 0 {
		return nil, core.ErrFileEmpty
	}

	if info.Size() > VALUE_MAX_SIZE {
		return nil, fmt.Errorf("file size exceeds maximum size %d", VALUE_MAX_SIZE)
	}
	return file, nil
}

// func AppendFromFile(name string, filePath string) error {
// 	// get strem by name hash
// 	// streamId := streamIdByName(streamName)
// 	iter := kvClientForPut.NewIterator(FILE_STREAM)
// 	// revert if exists
// 	iter.SeekToLast()
// 	if !iter.Valid() {
// 		return errors.New("The name unexists")
// 	}

// 	lastKey := string(iter.KeyValue().Key)
// 	lastFileId, err := parseStreamKey(lastKey)
// 	if err != nil {
// 		return err
// 	}

// 	// otherwise upload file
// 	if err := uploadFile(filePath, nil); err != nil {
// 		return errors.New("Failed to upload file")
// 	}
// 	// write stream with key0 be file root hash

// 	rootHash, err := GetRootHash(filePath)
// 	if err != nil {
// 		return errors.WithMessage(err, "Failed to get root hash")
// 	}

// 	batcher := kvClientForPut.Batcher()
// 	batcher.Set(FILE_STREAM,
// 		[]byte(getStreamKey(lastFileId+1)),
// 		[]byte(rootHash[:]),
// 	)

// 	err = batcher.Exec()
// 	if err != nil {
// 		return errors.WithMessage(err, "Failed to execute batcher")
// 	}
// 	return nil
// }

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
