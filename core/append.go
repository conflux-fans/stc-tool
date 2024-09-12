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

func AppendData(account common.Address, name string, data string) error {
	return appendDataOrCreate(account, name, data, false)
}

// appendDataOrCreate 向现有内容追加数据或在指定时创建新内容
//
// 参数:
//   - account: 操作账户地址
//   - name: 内容名称
//   - data: 要追加的数据
//   - createIfNotExist: 如果内容不存在是否创建新内容
func appendDataOrCreate(account common.Address, name string, data string, createIfNotExist bool) error {
	if len(data) > constants.VALUE_MAX_SIZE {
		return errors.New("Exceed max size once uploadable")
	}
	logger.Get().WithField("name", name).Info("Start append content")

	// split content to chunks
	var chunks []string
	for i := 0; i < len(data); i += constants.CHUNK_SIZE {
		end := lo.Min([]int{(i + 1) * constants.CHUNK_SIZE, len(data)})
		chunks = append(chunks, data[i*constants.CHUNK_SIZE:end])
	}

	meta, err := GetContentMetadata(name)
	if err != nil {
		if !(err == ERR_UNEXIST_CONTENT && createIfNotExist) {
			return err
		}
	}
	if meta == nil {
		meta = &ContentMetadata{
			LineSizeKey: keyLineCount(name),
			LineKeys:    []string{},
			LineSize:    0,
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
	adminBatcher := adminKvClientForPut.Batcher()
	batcher, err := getKvClientBatcher(account)
	if err != nil {
		return errors.WithMessage(err, "Failed to get kv client")
	}

	lineSizeKey := []byte(meta.LineSizeKey)
	lineSizeVal := []byte(fmt.Sprintf("%d", meta.LineSize+len(chunks)))

	batcher.Set(STREAM_FILE, lineSizeKey, lineSizeVal)
	adminBatcher.SetKeyToSpecial(STREAM_FILE, lineSizeKey).GrantSpecialWriteRole(STREAM_FILE, lineSizeKey, account)

	logger.Get().WithField("line size key", string(lineSizeKey)).Info("Set line size kv")

	for i, chunk := range chunks {
		key := []byte(keyLineIndex(name, meta.LineSize+i))

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

func AppendFromFile(account common.Address, name string, filePath string) error {
	return appendFromFileOrCreate(account, name, filePath, false)
}

// appendFromFileOrCreate 将文件内容追加到指定名称的目标中
// 参数:
//   - account: 账户地址
//   - name: 目标名称
//   - filePath: 要追加的文件路径
//   - force: 如果为true,则在目标不存在时创建新的目标并追加内容
func appendFromFileOrCreate(account common.Address, name string, filePath string, force bool) error {
	f, err := openFile(filePath)
	if err != nil {
		return err
	}

	// split by VALUE_MAX_SIZE
	for {
		buffer := make([]byte, constants.VALUE_MAX_SIZE)
		n, err := f.Read(buffer)
		if err != nil {
			return err
		}
		if n == 0 {
			return nil
		}

		if err = appendDataOrCreate(account, name, string(buffer[:n]), force); err != nil {
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

	if info.Size() > constants.VALUE_MAX_SIZE {
		return nil, fmt.Errorf("file size exceeds maximum size %d", constants.VALUE_MAX_SIZE)
	}
	return file, nil
}

func keyLineCount(name string) string {
	return fmt.Sprintf("%s:line", name)
}

func keyLineIndex(name string, index int) string {
	return fmt.Sprintf("%s:%d", name, index)
}
