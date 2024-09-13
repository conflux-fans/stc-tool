package core

import (
	"fmt"
	"os"

	"github.com/0glabs/0g-storage-client/core"
	ccore "github.com/0glabs/0g-storage-client/core"
	"github.com/conflux-fans/storage-cli/constants"
	"github.com/conflux-fans/storage-cli/constants/enums"
	"github.com/conflux-fans/storage-cli/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

var appender Appender

type Appender struct{}

func DefaultAppender() *Appender {
	return &appender
}

func (a *Appender) splitStringToChunks(str string) []string {
	var chunks []string
	for i := 0; i < len(str); i += constants.CHUNK_SIZE {
		end := lo.Min([]int{(i + 1) * constants.CHUNK_SIZE, len(str)})
		chunks = append(chunks, str[i*constants.CHUNK_SIZE:end])
	}
	return chunks
}

func (a *Appender) splitDataToChunks(data ccore.IterableData) []string {
	// split content to chunks
	var chunks []string
	iterator := data.Iterate(0, int64(constants.CONTENT_MAX_SIZE), false)
	for {
		chunk := iterator.Current()
		chunks = append(chunks, string(chunk))
	}
	return chunks
}

// appendExtendOrCreate 向现有内容追加数据或在指定时创建新内容
//
// 参数:
//   - account: 操作账户地址
//   - name: 内容名称
//   - data: 要追加的数据
//   - createIfNotExist: 如果内容不存在是否创建新内容
func (a *Appender) appendExtendOrCreate(account common.Address, name string, dataType enums.ExtendDataType, data ccore.IterableData, createIfNotExist bool) error {
	if data.Size() > constants.CONTENT_MAX_SIZE {
		return errors.New("Exceed max size once uploadable")
	}

	logger.Get().WithField("name", name).Info("Start append content")

	meta, err := a.GenMeta(account, name, createIfNotExist)
	if err != nil {
		return err
	}

	switch dataType {
	case enums.EXTEND_DATA_TEXT:
		return a.uploadExtendAsValue(account, name, meta, data)
	case enums.EXTEND_DATA_POINTER:
		return a.uploadExtendAsPointer(account, name, meta, data)
	}
	return nil
}

func (a *Appender) GenMeta(account common.Address, name string, createIfNotExist bool) (*ContentMetadata, error) {
	meta, err := GetContentMetadata(name)
	if err != nil {
		if err == ERR_UNEXIST_CONTENT && createIfNotExist {
			meta = &ContentMetadata{
				Name: name,
			}
		} else {
			return nil, err
		}
	} else {
		if !createIfNotExist {
			isWriter, err := CheckIsContentWriter(name, account)
			if err != nil {
				return nil, err
			}
			if !isWriter {
				return nil, fmt.Errorf("account %v is not writer of content", account)
			}
		}
		// 如果是 pointer 类型，报错
		if meta.ExtendDataType == enums.EXTEND_DATA_POINTER {
			return nil, errors.New("pointer type content not support append")
		}
	}

	return meta, nil
}

func (a *Appender) uploadExtendAsPointer(account common.Address, name string, meta *ContentMetadata, data ccore.IterableData) error {
	// 首先将整个文件上传
	mt, err := DefaultUploader().UploadIteratorData(data)
	if err != nil {
		return err
	}

	hashData, err := ccore.NewDataInMemory(mt.Root().Bytes())
	if err != nil {
		return err
	}
	// 将文件hash作为数据上传
	err = a.uploadExtendAsValue(account, name, meta, hashData)
	if err != nil {
		return errors.WithMessage(err, "上传文件hash失败")
	}
	return nil
}

func (a *Appender) uploadExtendAsValue(account common.Address, name string, meta *ContentMetadata, data ccore.IterableData) error {
	batcher, err := getKvClientBatcher(account)
	if err != nil {
		return errors.WithMessage(err, "Failed to get kv client")
	}

	iterator := data.Iterate(0, int64(constants.CHUNK_SIZE), false)

	i := 0
	for {
		exist, err := iterator.Next()
		if err != nil {
			return errors.WithMessage(err, "迭代数据时出错")
		}
		if !exist {
			break
		}
		chunk := iterator.Current()
		key := []byte(meta.LineIndexKey(meta.LineTotal + i))
		batcher.Set(STREAM_FILE, key, chunk)
		logger.Get().WithField("key", string(key)).Info("Set line kv")
		i++
	}
	logger.Get().WithField("name", name).Info("Set content values")

	lineTotalVal := []byte(fmt.Sprintf("%d", meta.LineTotal+i))
	batcher.Set(STREAM_FILE, []byte(meta.LineTotalKey()), lineTotalVal)
	batcher.Set(STREAM_FILE, []byte(meta.ExtendDataTypeKey()), []byte(meta.ExtendDataType.String()))
	batcher.Set(STREAM_FILE, []byte(meta.ExtendDataOwnerTokenIDKey()), []byte(fmt.Sprintf("%d", meta.OwnerTokenID)))
	logger.Get().Info("Set line metadata kvs")

	err = batcher.Exec()
	if err != nil {
		return errors.WithMessage(err, "Failed to set values of content")
	}

	logger.Get().WithField("name", name).WithField("line", i).Info("Append content completed")

	return nil
}

func (a *Appender) uploadLines(account common.Address, name string, meta *ContentMetadata, chunks []string) error {

	batcher, err := getKvClientBatcher(account)
	if err != nil {
		return errors.WithMessage(err, "Failed to get kv client")
	}

	lineCountVal := []byte(fmt.Sprintf("%d", meta.LineTotal+len(chunks)))

	batcher.Set(STREAM_FILE, []byte(meta.LineTotalKey()), lineCountVal)
	logger.Get().WithField("line count key", string(meta.LineTotalKey())).Info("Set line count kv")

	for i, chunk := range chunks {
		key := []byte(meta.LineIndexKey(meta.LineTotal + i))
		batcher.Set(STREAM_FILE, key, []byte(chunk))
		logger.Get().WithField("key", string(key)).Info("Set line kv")
	}

	logger.Get().WithField("name", name).Info("Set content values")
	err = batcher.Exec()
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
	adminBatcher := adminKvClientForPut.Batcher()
	batcher, err := getKvClientBatcher(account)
	if err != nil {
		return errors.WithMessage(err, "Failed to get kv client")
	}

	lineTotalKey := []byte(meta.LineTotalKey())
	lineTotalVal := []byte(fmt.Sprintf("%d", meta.LineTotal+len(chunks)))

	batcher.Set(STREAM_FILE, lineTotalKey, lineTotalVal)
	adminBatcher.SetKeyToSpecial(STREAM_FILE, lineTotalKey).GrantSpecialWriteRole(STREAM_FILE, lineTotalKey, account)

	logger.Get().WithField("line count key", string(lineTotalKey)).Info("Set line count kv")

	for i, chunk := range chunks {
		key := []byte(meta.LineIndexKey(meta.LineTotal + i))

		batcher.Set(STREAM_FILE, key, []byte(chunk))
		adminBatcher.SetKeyToSpecial(STREAM_FILE, []byte(key)).GrantSpecialWriteRole(STREAM_FILE, []byte(key), account)

		logger.Get().WithField("key", string(key)).Info("Set line kv")
	}

	logger.Get().WithField("name", name).Info("Set content keys to special keys")
	err = adminBatcher.Exec()
	if err != nil {
		return errors.WithMessage(err, "Failed to set speicial key by admin batcher")
	}

	logger.Get().WithField("name", name).Info("Set content values")
	err = batcher.Exec()
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
