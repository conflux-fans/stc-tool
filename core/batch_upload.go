package core

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/0glabs/0g-storage-client/core"
	"github.com/0glabs/0g-storage-client/kv"
	"github.com/0glabs/0g-storage-client/transfer"
	"github.com/conflux-fans/storage-cli/constants/enums"
	"github.com/conflux-fans/storage-cli/encrypt"
	"github.com/conflux-fans/storage-cli/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

const (
	ONE_BATCH_COUNT = 30000
)

type EncryptOption struct {
	Method   enums.CipherMethod
	Password string
}

func NewEncryptOption(method enums.CipherMethod, password string) (*EncryptOption, error) {
	if method == enums.CIPHER_METHOD_EMPTY && password == "" {
		return nil, nil
	}
	if method == enums.CIPHER_METHOD_EMPTY {
		return nil, errors.New("Missing cipher method specified")
	}
	if password == "" {
		return nil, errors.New("Missing password specified")
	}
	return &EncryptOption{
		Method:   method,
		Password: password,
	}, nil
}

// logger.Get().WithField("name", name).WithField("time use", timeUse).WithField("tps", count/int(timeUse/time.Second)).WithField("count", count).Info("Batch upload completed")
type BatchUploadResult struct {
	Name    string
	Count   int
	UseTime time.Duration
	TPS     int
}

func BatchUploadByKv(count int) error {
	GrantAllAccountStreamWriter()

	limit := len(kvBatcherForPut) * ONE_BATCH_COUNT
	if count > limit {
		return fmt.Errorf("exceed limit, the max limit batch upload count is %d", limit)
	}
	// name be time, gen 100000 kv
	name := fmt.Sprintf("BATCH-TEST-%d", time.Now().Unix())
	batchers := []*kv.Batcher{}

	// execute, every segment 20000 kv\
	m := ContentMetadata{
		Name:           name,
		ExtendDataType: enums.EXTEND_DATA_TEXT,
		LineTotal:      count,
	}
	for i := 0; i < count; {
		account := accounts[i/ONE_BATCH_COUNT]
		batcher := kvBatcherForPut[account]
		if i == 0 {
			batcher.Set(kvStreamId, []byte(m.LineTotalKey()), []byte(fmt.Sprintf("%d", m.LineTotal)))
		}

		// check account is writer, panic if not
		isWriter, err := CheckIsStreamWriter(account)
		if err != nil {
			return err
		}
		if !isWriter {
			return fmt.Errorf("account %s is not stream writer", account)
		}

		end := lo.Min([]int{count, i + ONE_BATCH_COUNT})
		for j := i; j < end; j++ {
			k, v := []byte(m.LineIndexKey(j)), []byte(fmt.Sprintf("%d", j))
			batcher.Set(kvStreamId, k, v)
			logger.Get().WithField("key", string(k)).WithField("value", string(v)).Debug("set key")
		}
		batchers = append(batchers, batcher)
		i = end
	}

	start := time.Now()
	logger.Get().WithField("time", start).WithField("name", name).WithField("len", len(batchers)).Info("Generated Datas")

	var w sync.WaitGroup
	var errs []error
	for _, b := range batchers {
		w.Add(1)
		go func(_b *kv.Batcher) {
			defer w.Done()
			_, err := _b.Exec(context.Background())
			if err != nil {
				errs = append(errs, err)
			}
		}(b)
	}
	w.Wait()

	var result BatchUploadResult
	if len(errs) > 0 {
		// logger.Get().WithError(errs[0]).WithField("count", count).Info("Failed to batch upload")
		return errors.WithMessage(errs[0], "Failed to batch upload")
	} else {
		timeUse := time.Since(start)
		tps := count / int(timeUse/time.Second)
		logger.Get().WithField("name", name).WithField("time use", timeUse).WithField("tps", tps).WithField("count", count).Info("Batch upload completed")

		result = BatchUploadResult{
			Name:    name,
			UseTime: timeUse,
			TPS:     tps,
			Count:   count,
		}
	}

	// query last
	fmt.Print("\x1b[36mINFO\x1b[0m[0000] \x1b[42m[TOOL]\x1b[0m Start verify ...")
	for i := 0; i < 1000; i++ {
		fmt.Print(".")
		v, err := kvClientForIterator.GetValue(context.Background(), kvStreamId, []byte(m.LineIndexKey(count-1)))
		if err != nil {
			logger.Get().WithError(err).Info("Failed to check upload state")
			time.Sleep(time.Millisecond * 100)
			continue
		}
		if v.Size == 0 {
			time.Sleep(time.Millisecond * 100)
			continue
		}

		fmt.Print("\n")
		logger.Get().WithField("last value", string(v.Data)).Info("Batch upload verified")
		logger.SuccessfWithParams(map[string]string{
			"Name":     result.Name,
			"Duration": result.UseTime.String(),
			"TPS":      fmt.Sprintf("%d", result.TPS),
			"Count":    fmt.Sprintf("%d", result.Count),
		}, "Batch upload completed and verified")
		return nil
	}
	logger.Fail("Failed to verify in 100 seconds")
	return nil
}

// TODO: count replace by source path?
func BatchUpload(count int, encryptOpt *EncryptOption) (common.Hash, error) {

	uploader, err := transfer.NewUploader(context.Background(), adminW3Client, zgNodeClients, zgLogOpt)
	if err != nil {
		return common.Hash{}, errors.WithMessage(err, "Failed to create uploader")
	}

	var datas []core.IterableData
	var report BatchUploadReport
	for i := 0; i < count; i++ {
		text, err := randomText(encryptOpt)
		if err != nil {
			return common.Hash{}, err
		}

		data, _ := core.NewDataInMemory(text)
		datas = append(datas, data)

		tree, err := core.MerkleTree(data)
		if err != nil {
			return common.Hash{}, errors.WithMessage(err, "Failed to create merkle tree")
		}
		report.Records = append(report.Records, uploadRecord{string(text), tree.Root()})
	}

	// save text and root to file
	report.StartTime = time.Now()
	hash, _, err := uploader.BatchUpload(context.Background(), datas, false)
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

	if encryptOpt == nil || encryptOpt.Method == enums.CIPHER_METHOD_EMPTY {
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
