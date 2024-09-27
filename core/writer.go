package core

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/conflux-fans/storage-cli/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

func CheckIsStreamWriter(account common.Address) (bool, error) {
	isWriter, err := kvClientForIterator.IsWriterOfStream(context.Background(), account, kvStreamId)
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

	isWriter := true

	keys := meta.AllKeys()
	var w sync.WaitGroup
	for i, lk := range keys {
		w.Add(1)
		go func(_lk []byte) {
			defer w.Done()

			_isWriter, err := kvClientForIterator.IsWriterOfKey(context.Background(), account, kvStreamId, _lk)
			if err != nil {
				panic(err)
			}
			if !_isWriter {
				logger.Get().WithField("key", string(_lk)).Info("Account is not writer of key")
				isWriter = false
			}
		}([]byte(lk))

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
	batcher := adminBatcher
	for _, account := range accounts {
		batcher.GrantWriteRole(kvStreamId, account)
	}

	_, err := batcher.Exec(context.Background())
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

	keys := meta.AllKeys()
	batcher := kvBatcherForPut[from]

	for _, k := range keys {
		batcher.GrantSpecialWriteRole(kvStreamId, []byte(k), to)
		batcher.RenounceSpecialWriteRole(kvStreamId, []byte(k))
	}

	_, err = batcher.Exec(context.Background())
	return err
}
