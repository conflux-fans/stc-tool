package core

import (
	"fmt"
	"os"

	ccore "github.com/0glabs/0g-storage-client/core"
	"github.com/0glabs/0g-storage-client/core/merkle"
	"github.com/0glabs/0g-storage-client/transfer"
	"github.com/conflux-fans/storage-cli/constants/enums"
	"github.com/conflux-fans/storage-cli/logger"
	"github.com/conflux-fans/storage-cli/utils/encryptutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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

var uploader Uploader

type Uploader struct{}

func DefaultUploader() *Uploader {
	return &uploader
}

// Upload data
func (u *Uploader) UploadByName(account common.Address, name string, dataType enums.ExtendDataType, data ccore.IterableData) error {
	logger.Get().WithField("name", name).Info("Ready to upload data")

	// revert if exists
	if err := u.checkDataNameExists(name); err != nil {
		return err
	}

	if err := DefaultAppender().appendExtendOrCreate(account, name, dataType, data, true); err != nil {
		return err
	}

	logger.Get().WithField("name", name).Info("Upload data completed")
	return nil
}

// func (u *Uploader) UploadDataFromFile(account common.Address, name string, filePath string) error {
// 	if err := u.checkDataNameExists(name); err != nil {
// 		return err
// 	}

// 	if err := DefaultAppender().appendDataFromFileOrCreate(account, name, filePath, true); err != nil {
// 		return err
// 	}

// 	logger.Get().WithField("name", name).Info("Upload data completed")
// 	return nil
// }

func (u *Uploader) checkDataNameExists(name string) error {
	// revert if exists
	m := ContentMetadata{
		Name: name,
	}

	v, err := kvClientForIterator.GetValue(STREAM_FILE, []byte(m.LineTotalKey()))
	if err != nil {
		return errors.WithMessage(err, "Failed to get file line size")
	}
	if v.Size != 0 {
		return fmt.Errorf("the name %s already exists", name)
	}
	return nil
}

func (u *Uploader) UploadFile(filepath string, opt *UploadOption) (*merkle.Tree, error) {
	if opt.EncryptOption != nil {
		outPath, err := encryptutils.EncryptFile(filepath, opt.EncryptOption.Method, opt.EncryptOption.Password)
		if err != nil {
			return nil, errors.WithMessage(err, "Failed to encrypt file")
		}
		filepath = outPath
		defer func() {
			os.Remove(outPath)
		}()
	}

	uploader := transfer.NewUploader(defaultFlow, nodeClients)

	f, err := ccore.Open(filepath)
	if err != nil {
		return nil, err
	}

	// Calculate file merkle root.
	tree, err := ccore.MerkleTree(f)
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to create data merkle tree")
	}
	logrus.WithField("root", tree.Root()).Info("Data merkle root calculated")

	if err = uploader.Upload(f); err != nil {
		return nil, err
	}

	return tree, nil
}

// upload data and return segments merkle tree and chunks merle tree
func (u *Uploader) UploadString(data []byte) (*merkle.Tree, *merkle.Tree, error) {
	uploader := transfer.NewUploader(defaultFlow, nodeClients)
	dataInMemory, err := ccore.NewDataInMemory(data)
	if err != nil {
		return nil, nil, err
	}

	err = uploader.Upload(dataInMemory)
	if err != nil && err.Error() != "Data already exists on ZeroGStorage network" {
		return nil, nil, err
	}

	segmentsTree, err := ccore.MerkleTree(dataInMemory)
	if err != nil {
		return nil, nil, err
	}

	chunksTree, err := u.getChunksTree(data)
	if err != nil {
		return nil, nil, err
	}

	return segmentsTree, chunksTree, nil
}

func (u *Uploader) UploadIteratorData(data ccore.IterableData) (*merkle.Tree, error) {
	uploader := transfer.NewUploader(defaultFlow, nodeClients)
	err := uploader.Upload(data)
	if err != nil {
		return nil, err
	}

	segmentsTree, err := ccore.MerkleTree(data)
	if err != nil {
		return nil, err
	}

	return segmentsTree, nil
}

func (u *Uploader) getChunksTree(data []byte) (*merkle.Tree, error) {
	dataInMemory, err := ccore.NewDataInMemory(data)
	if err != nil {
		return nil, err
	}

	batch := ccore.DefaultSegmentSize
	buf, err := ccore.ReadAt(dataInMemory, batch, 0, dataInMemory.PaddedSize())
	if err != nil {
		return nil, err
	}
	return u.buildChunksTree(buf), nil
}

func (u *Uploader) buildChunksTree(chunks []byte, emptyChunksPadded ...uint64) *merkle.Tree {
	var builder merkle.TreeBuilder

	// append chunks
	for offset, dataLen := 0, len(chunks); offset < dataLen; offset += ccore.DefaultChunkSize {
		chunk := chunks[offset : offset+ccore.DefaultChunkSize]
		builder.Append(chunk)
	}

	// append empty chunks
	if len(emptyChunksPadded) > 0 && emptyChunksPadded[0] > 0 {
		for i := uint64(0); i < emptyChunksPadded[0]; i++ {
			builder.AppendHash(ccore.EmptyChunkHash)
		}
	}

	tree := builder.Build()
	return tree
}
