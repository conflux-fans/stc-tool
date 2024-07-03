package core

import (
	"fmt"
	"sync"
	"time"

	"github.com/conflux-fans/storage-cli/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

func CheckIsStreamWriter(account common.Address) (bool, error) {
	isWriter, err := kvClientForIterator.IsWriterOfStream(account, STREAM_FILE)
	if err != nil {
		return false, err
	}
	return isWriter, nil
}

func CheckIsContentWriter(name string, account common.Address) (bool, error) {
	meta, err := GetContentMetadata(name)
	if err != nil {
		return false, err
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
	return isWriter, nil
}

func GrantStreamWriter(accounts ...common.Address) error {
	allAreWriter := true
	for _, account := range accounts {
		isWriter, err := CheckIsStreamWriter(account)
		if err != nil {
			return errors.WithMessage(err, "failed to check stream writer")
		}
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

	err := batcher.Exec()
	if err != nil {
		return err
	}

	if err = waitStreamWriterConfirm(accounts[len(accounts)-1]); err != nil {
		return err
	}

	logger.Get().Info("Grant writers done")

	return nil
}

func waitStreamWriterConfirm(account common.Address) error {
	logger.Get().Info("Wait write setting")
	for i := 0; i < 100; i++ {
		time.Sleep(time.Second)
		isWriter, err := CheckIsStreamWriter(account)
		if isWriter {
			return nil
		}

		if i == 100-1 {
			if err != nil {
				return errors.WithMessage(err, "failed to grant stream writer")
			}
		}
	}
	return errors.New("failed to grant stream writer")
}

func TransferWriter(name string, from common.Address, to common.Address) error {
	// get all keys
	logger.Get().WithField("name", name).WithField("from", from).WithField("to", to).Info("Start transfer content owner")
	isWriter, err := CheckIsContentWriter(name, from)
	if err != nil {
		return errors.WithMessage(err, "failed to check if content owner")
	}

	if !isWriter {
		return fmt.Errorf("account %v is not the writer of content %v", from, name)
	}

	if from == to {
		return fmt.Errorf("from account should not same with to account")
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
