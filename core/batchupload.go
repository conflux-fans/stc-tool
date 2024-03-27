package core

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/zero-gravity-labs/zerog-storage-client/core"
	"github.com/zero-gravity-labs/zerog-storage-client/kv"
	"github.com/zero-gravity-labs/zerog-storage-client/transfer"
	"github.com/zero-gravity-labs/zerog-storage-tool/encrypt"
)

type EncryptOption struct {
	Method   string
	Password string
}

func NewEncryptOption(method string, password string) (*EncryptOption, error) {
	if method == "" && password == "" {
		return nil, nil
	}
	if method == "" {
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

func BatchUploadByKv(count int) {
	// name be time, gen 100000 kv
	name := fmt.Sprintf("%d", time.Now().Unix())

	batchers := []*kv.Batcher{}

	// execute, every segment 1000 kv\
	for i := 0; i < count; {
		batcher := defaultKvClientForPut.Batcher()

		end := lo.Min([]int{count, i + 100000})
		for j := i; j < end; j++ {
			batcher.Set(STREAM_FILE,
				[]byte(keyLineIndex(name, j)),
				[]byte(fmt.Sprintf("%d", j)),
			)
		}
		batchers = append(batchers, batcher)
		i = end
	}

	logrus.WithField("len", len(batchers)).Info("Generated Datas")
	start := time.Now()

	var w sync.WaitGroup
	var errs []error
	for _, b := range batchers {
		w.Add(1)
		go func(_b *kv.Batcher) {
			defer w.Done()
			err := _b.Exec()
			if err != nil {
				errs = append(errs, err)
			}
		}(b)
	}
	w.Wait()
	if len(errs) > 0 {
		logrus.WithError(errs[0]).WithField("count", count).Info("Failed to batch upload")
		return
	} else {
		timeUse := time.Since(start)
		logrus.WithField("time use", timeUse).WithField("count", count).Info("Batch upload completed")
	}

	// query last
	for i := 0; i < 1000; i++ {
		fmt.Print(".")
		v, err := kvClientForIterator.GetValue(STREAM_FILE, []byte(keyLineIndex(name, count-1)))
		if err != nil {
			logrus.WithError(err).Info("Failed to check upload state")
			time.Sleep(time.Millisecond * 100)
			continue
		}
		if v.Size == 0 {
			time.Sleep(time.Millisecond * 100)
			continue
		}
		fmt.Print("\n")
		logrus.WithField("value", string(v.Data)).Info("Batch upload verified")
		break
	}
}

// TODO: count replace by source path?
func BatchUpload(count int, encryptOpt *EncryptOption) (common.Hash, error) {

	uploader := transfer.NewUploader(defaultFlow, nodeClients)

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

	if encryptOpt == nil || encryptOpt.Method == "" {
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
