package core

import (
	"context"
	"fmt"
	"os"

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
		return nil, ccore.ErrFileRequired
	}

	if info.Size() == 0 {
		return nil, ccore.ErrFileEmpty
	}

	if info.Size() > constants.CONTENT_MAX_SIZE {
		return nil, fmt.Errorf("file size exceeds maximum size %d", constants.CONTENT_MAX_SIZE)
	}
	return file, nil
}
