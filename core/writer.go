package core

import (
	"fmt"
	"sync"

	"github.com/conflux-fans/storage-cli/logger"
	"github.com/ethereum/go-ethereum/common"
)

func CheckIsStreamWriter(account common.Address) bool {
	isWriter, err := kvClientForIterator.IsWriterOfStream(account, STREAM_FILE)
	if err != nil {
		panic(err)
	}
	return isWriter
}

func CheckIsContentWriter(name string, account common.Address) bool {
	meta, err := GetContentMetadata(name)
	if err != nil {
		panic(err)
	}

	keys := append(meta.LineKeys, meta.LineSizeKey)
	isWriter := true

	var w sync.WaitGroup
	for i, lk := range keys {
		w.Add(1)
		go func(_lk string) {
			defer w.Done()

			_isWriter, err := kvClientForIterator.IsWriterOfKey(account, STREAM_FILE, []byte(_lk))
			if err != nil {
				panic(err)
			}
			if !_isWriter {
				logger.Get().WithField("key", string(_lk)).Info("Account is not writer of key")
				isWriter = false
			}
		}(lk)

		if i%20 == 0 || i == len(keys)-1 {
			w.Wait()
			if !isWriter {
				break
			}
		}
	}
	return isWriter
}

func GrantStreamWriter(accounts ...common.Address) error {
	allAreWriter := true
	for _, account := range accounts {
		isWriter := CheckIsStreamWriter(account)
		allAreWriter = isWriter && allAreWriter
	}
	if allAreWriter {
		logger.Get().Info("All accounts are writer of stream")
		return nil
	}

	logger.Get().WithField("accounts", accounts).Info("Grant stream writer to accounts")
	batcher := adminKvClientForPut.Batcher()
	for _, account := range accounts {
		batcher.GrantWriteRole(STREAM_FILE, account)
	}
	logger.Get().Info("Grant writers done")
	return batcher.Exec()
}

func TransferWriter(name string, from common.Address, to common.Address) error {
	// get all keys
	logger.Get().WithField("name", name).WithField("from", from).WithField("to", to).Info("Start transfer content owner")
	isWriter := CheckIsContentWriter(name, from)
	if !isWriter {
		return fmt.Errorf("account %v not the writer of content %v", from, name)
	}

	meta, err := GetContentMetadata(name)
	if err != nil {
		return err
	}

	keys := append(meta.LineKeys, meta.LineSizeKey)
	batcher := kvClientsForPut[from].Batcher()

	for _, k := range keys {
		batcher.GrantSpecialWriteRole(STREAM_FILE, []byte(k), to)
		batcher.RenounceSpecialWriteRole(STREAM_FILE, []byte(k))
	}

	return batcher.Exec()
}
