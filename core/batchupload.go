package core

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/zero-gravity-labs/zerog-storage-client/core"
	"github.com/zero-gravity-labs/zerog-storage-client/transfer"
	"github.com/zero-gravity-labs/zerog-storage-tool/encrypt"
)

type EncryptOption struct {
	Method   string
	Password string
}

// TODO: count replace by source path?
func BatchUpload(count int, encryptOpt *EncryptOption) (common.Hash, error) {

	uploader := transfer.NewUploader(flow, nodeClients)

	var datas []core.IterableData
	var report BatchUploadReport
	for i := 0; i < count; i++ {
		text, err := randomText(encryptOpt)
		if err != nil {
			return common.Hash{}, err
		}

		data := core.NewDataInMemory(text)
		datas = append(datas, data)

		tree, err := core.MerkleTree(data)
		if err != nil {
			return common.Hash{}, errors.WithMessage(err, "Failed to create merkle tree")
		}
		report.Records = append(report.Records, uploadRecord{string(text), tree.Root()})
	}

	// save text and root to file
	report.StartTime = time.Now()
	hash, err := uploader.BatchUpload(datas, false)
	if err != nil {
		return common.Hash{}, errors.WithMessage(err, "Failed to batch upload")
	}
	report.EndTime = time.Now()
	report.Hash = hash

	if err := report.Save(fmt.Sprintf("./log/%s.json", time.Now().Format(time.RFC3339))); err != nil {
		return common.Hash{}, errors.WithMessage(err, "Failed to save report")
	}
	return hash, nil
}

func randomText(encryptOpt *EncryptOption) ([]byte, error) {
	rnd := rand.Intn(10000)
	content := []byte(fmt.Sprintf("%v - %d - hello world", time.Now().Format(time.RFC3339Nano), rnd))

	if encryptOpt.Method == "" {
		return content, nil
	}

	encryptor, err := encrypt.GetEncryptor(encryptOpt.Method)
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to get encryptor")
	}

	ciphertext, err := encrypt.EncryptBytes(encryptor, content, []byte(encryptOpt.Password))
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to encrypt")
	}

	return ciphertext, nil
}
