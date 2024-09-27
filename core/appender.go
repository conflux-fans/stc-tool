package core

import (
	"context"
	"fmt"
	"os"

	"github.com/0glabs/0g-storage-client/core"
	ccore "github.com/0glabs/0g-storage-client/core"
	"github.com/conflux-fans/storage-cli/constants"
	"github.com/conflux-fans/storage-cli/constants/enums"
	"github.com/conflux-fans/storage-cli/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

var appender Appender

type Appender struct{}

func DefaultAppender() *Appender {
	return &appender
}

// appendExtendOrCreate 向现有内容追加数据或在指定时创建新内容
//
// 参数:
//   - account: 操作账户地址
//   - name: 内容名称
//   - data: 要追加的数据
//   - createIfNotExist: 如果内容不存在是否创建新内容
// func (a *Appender) appendExtendOrCreate(account common.Address, name string, dataType enums.ExtendDataType, data ccore.IterableData, createIfNotExist bool) error {
// 	if data.Size() > constants.CONTENT_MAX_SIZE {
// 		return errors.New("Exceed max size once uploadable")
// 	}

// 	logger.Get().WithField("name", name).Info("Start append content")

// 	meta, err := a.GetMeta(account, name)
// 	if err != nil {
// 		if errors.Cause(err) == ERR_UNEXIST_CONTENT && createIfNotExist {
// 			txHash, tokenID, err := DefaultOwnerOperator().Mint(account)
// 			if err != nil {
// 				return errors.WithMessage(err, "Failed to mint")
// 			}
// 			logger.Get().WithField("txHash", txHash.Hex()).WithField("tokenID", tokenID.String()).Info("Mint content owner NFT completed")

// 			meta = &ContentMetadata{
// 				Name:           name,
// 				ExtendDataType: dataType,
// 				OwnerTokenID:   tokenID.Uint64(),
// 			}
// 		} else {
// 			return err
// 		}
// 	}

// 	switch dataType {
// 	case enums.EXTEND_DATA_TEXT:
// 		return a.uploadExtendAsText(account, name, meta, data)
// 	case enums.EXTEND_DATA_POINTER:
// 		return a.uploadExtendAsPointer(account, name, meta, data)
// 	}
// 	return nil
// }

func (a *Appender) AppendExtend(account common.Address, name string, data ccore.IterableData) error {
	if data.Size() > constants.CONTENT_MAX_SIZE {
		return errors.New("Exceed max size once uploadable")
	}

	logger.Get().WithField("name", name).Info("Start append content")

	meta, err := a.GetMeta(account, name)
	if err != nil {
		return err
	}

	return DefaultUploader().uploadExtend(account, name, meta, data)
}

func (a *Appender) GetMeta(account common.Address, name string) (*ContentMetadata, error) {
	meta, err := GetContentMetadata(name)
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to get content metadata")
	}

	// if err := a.checkWritePermission(name, account); err != nil {
	// 	return nil, err
	// }

	// Check if owner
	isOwner, err := DefaultOwnerOperator().CheckIsContentOwner(account, name)
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to check if owner")
	}
	if !isOwner {
		return nil, errors.New("Account is not the content owner")
	}

	if meta.ExtendDataType == enums.EXTEND_DATA_POINTER {
		return nil, errors.New("Pointer type content does not support append operation")
	}

	return meta, nil
}

func (a *Appender) checkWritePermission(name string, account common.Address) error {
	isWriter, err := CheckIsContentWriter(name, account)
	if err != nil {
		return fmt.Errorf("error checking write permission: %v", err)
	}
	if !isWriter {
		return fmt.Errorf("account %v is not a writer of the content", account)
	}
	return nil
}

func (a *Appender) uploadExtendAsText(account common.Address, name string, meta *ContentMetadata, data ccore.IterableData) error {
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
			return errors.WithMessage(err, "迭代数据时出错")
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

func (a *Appender) uploadExtendAsPointer(account common.Address, name string, meta *ContentMetadata, data ccore.IterableData) error {
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

func (a *Appender) uploadStringLines(account common.Address, name string, meta *ContentMetadata, chunks []string) error {

	batcher, err := getKvBatcher(account)
	if err != nil {
		return errors.WithMessage(err, "Failed to get kv client")
	}

	lineCountVal := []byte(fmt.Sprintf("%d", meta.LineTotal+len(chunks)))

	batcher.Set(kvStreamId, []byte(meta.LineTotalKey()), lineCountVal)
	logger.Get().WithField("line count key", string(meta.LineTotalKey())).Info("Set line count kv")

	for i, chunk := range chunks {
		key := []byte(meta.LineIndexKey(meta.LineTotal + i))
		batcher.Set(kvStreamId, key, []byte(chunk))
		logger.Get().WithField("key", string(key)).Info("Set line kv")
	}

	logger.Get().WithField("name", name).Info("Set content values")
	_, err = batcher.Exec(context.Background())
	if err != nil {
		return errors.WithMessage(err, "Failed to set values of content")
	}

	logger.Get().WithField("name", name).WithField("line", len(chunks)).Info("Append content completed")

	return nil
}

// uploadLinesAndSetSpecialWriter 批量上传内容并设置每行的 special writer
//
// 参数:
//   - account: 操作账户地址
//   - name: 内容名称
//   - meta: 内容元数据
//   - chunks: 要上传的数据块切片
//
// 返回:
//   - error: 如果上传过程中出现错误则返回相应的错误信息
func (a *Appender) uploadLinesAndSetSpecialWriter(account common.Address, name string, meta *ContentMetadata, chunks []string) error {
	// adminBatcher := kv.NewBatcher(math.MaxUint64, zgNodeClients) //adminKvClientForPut.Batcher()
	batcher, err := getKvBatcher(account)
	if err != nil {
		return errors.WithMessage(err, "Failed to get kv client")
	}

	lineTotalKey := []byte(meta.LineTotalKey())
	lineTotalVal := []byte(fmt.Sprintf("%d", meta.LineTotal+len(chunks)))

	batcher.Set(kvStreamId, lineTotalKey, lineTotalVal)
	adminBatcher.SetKeyToSpecial(kvStreamId, lineTotalKey).GrantSpecialWriteRole(kvStreamId, lineTotalKey, account)

	logger.Get().WithField("line count key", string(lineTotalKey)).Info("Set line count kv")

	for i, chunk := range chunks {
		key := []byte(meta.LineIndexKey(meta.LineTotal + i))

		batcher.Set(kvStreamId, key, []byte(chunk))
		adminBatcher.SetKeyToSpecial(kvStreamId, []byte(key)).GrantSpecialWriteRole(kvStreamId, []byte(key), account)

		logger.Get().WithField("key", string(key)).Info("Set line kv")
	}

	logger.Get().WithField("name", name).Info("Set content keys to special keys")
	_, err = adminBatcher.Exec(context.Background())
	if err != nil {
		return errors.WithMessage(err, "Failed to set speicial key by admin batcher")
	}

	logger.Get().WithField("name", name).Info("Set content values")
	_, err = batcher.Exec(context.Background())
	if err != nil {
		return errors.WithMessage(err, "Failed to set values of content")
	}

	logger.Get().WithField("name", name).WithField("line", len(chunks)).Info("Append content completed")

	return nil
}

// func (a *Appender) AppendDataFromFile(account common.Address, name string, filePath string) error {
// 	return a.appendDataFromFileOrCreate(account, name, filePath, false)
// }

// // appendDataFromFileOrCreate 将文件内容追加到指定名称的目标中
// // 参数:
// //   - account: 账户地址
// //   - name: 目标名称
// //   - filePath: 要追加的文件路径
// //   - force: 如果为true,则在目标不存在时创建新的目标并追加内容
// func (a *Appender) appendDataFromFileOrCreate(account common.Address, name string, filePath string, force bool) error {
// 	f, err := a.openFile(filePath)
// 	if err != nil {
// 		return err
// 	}

// 	// split by VALUE_MAX_SIZE
// 	for {
// 		buffer := make([]byte, constants.CONTENT_MAX_SIZE)
// 		n, err := f.Read(buffer)
// 		if err != nil {
// 			return err
// 		}
// 		if n == 0 {
// 			return nil
// 		}

// 		if err = a.appendDataOrCreate(account, name, string(buffer[:n]), force); err != nil {
// 			return err
// 		}
// 	}
// }

func (a *Appender) openFile(name string) (*os.File, error) {
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

	if info.Size() > constants.CONTENT_MAX_SIZE {
		return nil, fmt.Errorf("file size exceeds maximum size %d", constants.CONTENT_MAX_SIZE)
	}
	return file, nil
}
