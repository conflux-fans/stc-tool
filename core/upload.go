package core

import (
	"context"
	"fmt"
	"os"

	ccore "github.com/0glabs/0g-storage-client/core"
	"github.com/0glabs/0g-storage-client/core/merkle"
	"github.com/0glabs/0g-storage-client/transfer"
	"github.com/conflux-fans/storage-cli/constants"
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
// func (u *Uploader) UploadExtend(account common.Address, name string, dataType enums.ExtendDataType, data ccore.IterableData) error {
// 	logger.Get().WithField("name", name).Info("Ready to upload content")

// 	// revert if exists
// 	if err := u.checkExtendNameExists(name); err != nil {
// 		return err
// 	}

// 	if err := u.Create(account, name, dataType, data, true); err != nil {
// 		return err
// 	}

// 	logger.Get().WithField("name", name).Info("Upload data completed")
// 	return nil
// }

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

func (a *Uploader) UploadExtendIfNotExist(account common.Address, name string, dataType enums.ExtendDataType, data ccore.IterableData) error {
	if data.Size() > constants.CONTENT_MAX_SIZE {
		return errors.New("Exceed max size once uploadable")
	}

	logger.Get().WithField("name", name).Info("Start append content")

	_, err := GetContentMetadata(name)
	if err == nil {
		return errors.New("content already exists")
	} else if err != ERR_UNEXIST_CONTENT {
		return errors.WithMessage(err, "Failed to get content metadata")
	}

	txHash, tokenID, err := DefaultOwnerOperator().Mint(account)
	if err != nil {
		return errors.WithMessage(err, "Failed to mint")
	}
	logger.Get().WithField("txHash", txHash.Hex()).WithField("tokenID", tokenID.String()).Info("Mint content owner NFT completed")

	meta := &ContentMetadata{
		Name:           name,
		ExtendDataType: dataType,
		OwnerTokenID:   tokenID.String(),
	}

	return a.uploadExtend(account, name, meta, data)
}

func (a *Uploader) uploadExtend(account common.Address, name string, meta *ContentMetadata, data ccore.IterableData) error {
	switch meta.ExtendDataType {
	case enums.EXTEND_DATA_TEXT:
		return a.uploadExtendAsText(account, name, meta, data)
	case enums.EXTEND_DATA_POINTER:
		return a.uploadExtendAsPointer(account, name, meta, data)
	}
	return fmt.Errorf("unsupported extend data type %v", meta.ExtendDataType)
}

func (a *Uploader) uploadExtendAsText(account common.Address, name string, meta *ContentMetadata, data ccore.IterableData) error {
	batcher, err := getKvBatcher(account)
	if err != nil {
		return errors.WithMessage(err, "Failed to get kv client")
	}

	iterator := data.Iterate(0, int64(constants.CHUNK_SIZE), false)

	entries := make(map[string]string)

	i := 0
	for {
		exist, err := iterator.Next()
		if err != nil {
			return errors.WithMessage(err, "Error iterating data")
		}
		if !exist {
			break
		}
		entries[meta.LineIndexKey(meta.LineTotal+i)] = string(iterator.Current())
		i++
	}
	logger.Get().WithField("line num", i).Info("content lines")

	entries[meta.LineTotalKey()] = fmt.Sprintf("%d", meta.LineTotal+i)
	entries[meta.ExtendDataTypeKey()] = meta.ExtendDataType.String()
	entries[meta.ExtendDataOwnerTokenIDKey()] = meta.OwnerTokenID

	for k, v := range entries {
		batcher.Set(kvStreamId, []byte(k), []byte(v))
	}
	logger.Get().WithField("entries", entries).Info("Set line metadata kvs")

	_, err = batcher.Exec(context.Background())
	if err != nil {
		return errors.WithMessage(err, "Failed to set values of content")
	}
	logger.Get().WithField("name", name).WithField("line", i).Info("Append content completed")

	return nil
}

func (a *Uploader) uploadExtendAsPointer(account common.Address, name string, meta *ContentMetadata, data ccore.IterableData) error {
	// First, upload the entire file
	mt, err := DefaultUploader().UploadIteratorData(data)
	if err != nil {
		return err
	}

	hashData, err := ccore.NewDataInMemory(mt.Root().Bytes())
	if err != nil {
		return err
	}
	// Upload the file hash as data
	err = a.uploadExtendAsText(account, name, meta, hashData)
	if err != nil {
		return errors.WithMessage(err, "Failed to upload file hash")
	}
	return nil
}

func (u *Uploader) checkExtendNameExists(name string) error {
	// revert if exists
	m := ContentMetadata{
		Name: name,
	}

	v, err := kvClientForIterator.GetValue(context.Background(), kvStreamId, []byte(m.LineTotalKey()))
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

	uploader, err := transfer.NewUploader(context.Background(), adminW3Client, zgNodeClients)
	if err != nil {
		return nil, err
	}

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

	if _, err = uploader.Upload(context.Background(), f); err != nil {
		return nil, err
	}

	return tree, nil
}

// upload data and return segments merkle tree and chunks merle tree
func (u *Uploader) UploadString(data []byte) (*merkle.Tree, *merkle.Tree, error) {
	uploader, err := transfer.NewUploader(context.Background(), adminW3Client, zgNodeClients)
	if err != nil {
		return nil, nil, err
	}
	dataInMemory, err := ccore.NewDataInMemory(data)
	if err != nil {
		return nil, nil, err
	}

	_, err = uploader.Upload(context.Background(), dataInMemory)
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
	uploader, err := transfer.NewUploader(context.Background(), adminW3Client, zgNodeClients)
	if err != nil {
		return nil, err
	}
	_, err = uploader.Upload(context.Background(), data)
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
