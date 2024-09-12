package core

import (
	"fmt"
	"os"

	"github.com/0glabs/0g-storage-client/core"
	"github.com/conflux-fans/storage-cli/constants"
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

func (a *Appender) AppendDataFromContent(account common.Address, name string, data string) error {
	return a.appendDataOrCreate(account, name, data, false)
}

func (a *Appender) splitContentToChunks(data string) []string {
	// split content to chunks
	var chunks []string
	for i := 0; i < len(data); i += constants.CHUNK_SIZE {
		end := lo.Min([]int{(i + 1) * constants.CHUNK_SIZE, len(data)})
		chunks = append(chunks, data[i*constants.CHUNK_SIZE:end])
	}
	return chunks
}

// appendDataOrCreate 向现有内容追加数据或在指定时创建新内容
//
// 参数:
//   - account: 操作账户地址
//   - name: 内容名称
//   - data: 要追加的数据
//   - createIfNotExist: 如果内容不存在是否创建新内容
func (a *Appender) appendDataOrCreate(account common.Address, name string, data string, createIfNotExist bool) error {
	if len(data) > constants.CONTENT_MAX_SIZE {
		return errors.New("Exceed max size once uploadable")
	}
	logger.Get().WithField("name", name).Info("Start append content")

	chunks := a.splitContentToChunks(data)

	meta, err := GetContentMetadata(name)
	if err != nil {
		if !(err == ERR_UNEXIST_CONTENT && createIfNotExist) {
			return err
		} else {
			meta = &ContentMetadata{
				LineCountKey: keyLineCount(name),
				LineCount:    0,
			}
		}
	}

	// error if not writer
	if !createIfNotExist {
		isWriter, err := CheckIsContentWriter(name, account)
		if err != nil {
			return err
		}
		if !isWriter {
			return fmt.Errorf("account %v is not writer of content", account)
		}
	}

	return a.uploadLines(account, name, meta, chunks)
}

func (a *Appender) uploadLines(account common.Address, name string, meta *ContentMetadata, chunks []string) error {

	batcher, err := getKvClientBatcher(account)
	if err != nil {
		return errors.WithMessage(err, "Failed to get kv client")
	}

	lineCountKey := []byte(meta.LineCountKey)
	lineCountVal := []byte(fmt.Sprintf("%d", meta.LineCount+len(chunks)))

	batcher.Set(STREAM_FILE, lineCountKey, lineCountVal)
	logger.Get().WithField("line count key", string(lineCountKey)).Info("Set line count kv")

	for i, chunk := range chunks {
		key := []byte(keyLineIndex(name, meta.LineCount+i))
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

	lineCountKey := []byte(meta.LineCountKey)
	lineCountVal := []byte(fmt.Sprintf("%d", meta.LineCount+len(chunks)))

	batcher.Set(STREAM_FILE, lineCountKey, lineCountVal)
	adminBatcher.SetKeyToSpecial(STREAM_FILE, lineCountKey).GrantSpecialWriteRole(STREAM_FILE, lineCountKey, account)

	logger.Get().WithField("line count key", string(lineCountKey)).Info("Set line count kv")

	for i, chunk := range chunks {
		key := []byte(keyLineIndex(name, meta.LineCount+i))

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

func (a *Appender) AppendDataFromFile(account common.Address, name string, filePath string) error {
	return a.appendDataFromFileOrCreate(account, name, filePath, false)
}

// appendDataFromFileOrCreate 将文件内容追加到指定名称的目标中
// 参数:
//   - account: 账户地址
//   - name: 目标名称
//   - filePath: 要追加的文件路径
//   - force: 如果为true,则在目标不存在时创建新的目标并追加内容
func (a *Appender) appendDataFromFileOrCreate(account common.Address, name string, filePath string, force bool) error {
	f, err := a.openFile(filePath)
	if err != nil {
		return err
	}

	// split by VALUE_MAX_SIZE
	for {
		buffer := make([]byte, constants.CONTENT_MAX_SIZE)
		n, err := f.Read(buffer)
		if err != nil {
			return err
		}
		if n == 0 {
			return nil
		}

		if err = a.appendDataOrCreate(account, name, string(buffer[:n]), force); err != nil {
			return err
		}
	}
}

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
