package core

import (
	"fmt"
	"math/rand"
	"os"
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

const (
	ONE_BATCH_COUNT = 20000
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
	limit := len(kvClientsForPut) * ONE_BATCH_COUNT
	if count > limit {
		fmt.Printf("Exceed limit, the max limit batch upload count is %ds\n", count)
	}
	// name be time, gen 100000 kv
	name := fmt.Sprintf("BATCH-TEST-%d", time.Now().Unix())
	batchers := []*kv.Batcher{}

	// execute, every segment 1000 kv\
	// kvClientForPutList := lo.Values(kvClientsForPut)
	for i := 0; i < count; {
		account := accounts[i/ONE_BATCH_COUNT]
		batcher := kvClientsForPut[account].Batcher()
		if i == 0 {
			batcher.Set(STREAM_FILE, []byte(keyLineCount(name)), []byte(fmt.Sprintf("%d", count)))
		}

		// check account is writer, panic if not
		if isWriter := CheckIsStreamWriter(account); !isWriter {
			fmt.Printf("Failed to batch upload: Account %s is not stream writer\n", account)
			os.Exit(1)
		}

		end := lo.Min([]int{count, i + ONE_BATCH_COUNT})
		for j := i; j < end; j++ {
			k, v := []byte(keyLineIndex(name, j)), []byte(fmt.Sprintf("%d", j))
			batcher.Set(STREAM_FILE, k, v)
			logrus.WithField("key", string(k)).WithField("value", string(v)).Debug("set key")
		}
		batchers = append(batchers, batcher)
		i = end
	}

	start := time.Now()
	logrus.WithField("time", start).WithField("name", name).WithField("len", len(batchers)).Info("Generated Datas")

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
		logrus.WithField("name", name).WithField("time use", timeUse).WithField("tps", count/int(timeUse/time.Second)).WithField("count", count).Info("Batch upload completed")
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
		return
	}
	fmt.Print("\n")
	fmt.Println("Failed to verify in 100 seconds")
}

// func panicIfNotStreamWriter(account common.Address) {
// 	isWriter, err := kvClientForIterator.IsWriterOfStream(account, STREAM_FILE)
// 	if err != nil {
// 		panic(err)
// 	}
// 	if !isWriter {
// 		panic(fmt.Sprintf("account %v is not writer", account))
// 	}
// }

// func mustCheckIsContentWriter(account common.Address, name string) bool {
// 	meta, err := GetContentMetadata(name)
// 	if err != nil {
// 		panic(err)
// 	}
// 	isWriter, err := kvClientForIterator.IsWriterOfKey(account, STREAM_FILE, []byte(meta.LineSizeKey))
// 	if err != nil {
// 		panic(err)
// 	}
// 	if !isWriter {
// 		panic(fmt.Sprintf("%v is not writer of key %s", account, meta.LineSizeKey))
// 	}

// 	var w sync.WaitGroup
// 	for _, lk := range meta.LineKeys {
// 		go func(_lk string) {
// 			w.Add(1)
// 			defer w.Done()

// 			isWriter, err := kvClientForIterator.IsWriterOfKey(account, STREAM_FILE, []byte(lk))
// 			if err != nil {
// 				panic(err)
// 			}
// 			if !isWriter {
// 				panic(fmt.Sprintf("%v is not writer of key %s", account, _lk))
// 			}
// 		}(lk)
// 	}
// 	w.Wait()
// 	return true
// }

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
