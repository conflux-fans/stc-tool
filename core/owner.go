package core

import (
	"math/big"

	"github.com/conflux-fans/storage-cli/pkg/utils/bigutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

var ownerOperator OwnerOperator

type OwnerOperator struct {
}

func DefaultOwnerOperator() *OwnerOperator {
	return &ownerOperator
}

// func TransferOwner(name string, from common.Address, to common.Address) error {
// 	// get all keys
// 	logger.Get().WithField("name", name).WithField("from", from).WithField("to", to).Info("Start transfer content owner")

// 	meta, err := GetContentMetadata(name)
// 	if err != nil {
// 		return err
// 	}

// 	keys := append(meta.LineKeys, meta.LineSizeKey)

// 	logger.Get().WithField("length", len(keys)).Info("Get content related keys")

// 	// check is all writer, if not
// 	for _, k := range keys {
// 		isWriter, err := kvClientForIterator.IsWriterOfKey(defaultAccount, kvStreamId, []byte(k))
// 		if err != nil {
// 			return errors.WithMessage(err, "Failed to check if owner")
// 		}
// 		if !isWriter {
// 			return fmt.Errorf("not the writer of key %s", k)
// 		}
// 	}

// 	batcher := kvClientsForPut[from].Batcher()

// 	for _, k := range keys {
// 		batcher.GrantSpecialWriteRole(kvStreamId, []byte(k), to)
// 		batcher.RenounceSpecialWriteRole(kvStreamId, []byte(k))
// 	}

// 	return batcher.Exec()
// }

func (o *OwnerOperator) Mint(to common.Address) (common.Hash, *big.Int, error) {
	receipt, tokenID, err := DefaultPmContractHelper().Mint(to)
	if err != nil {
		return common.Hash{}, nil, errors.WithMessage(err, "Failed to invoke mint of contract")
	}

	return receipt.TransactionHash, tokenID, nil
}

func (o *OwnerOperator) TransferOwner(name string, from common.Address, to common.Address) (common.Hash, error) {
	meta, err := GetContentMetadata(name)
	if err != nil {
		return common.Hash{}, errors.WithMessage(err, "Failed to get content metadata")
	}

	tokenID := bigutils.MustParseBigInt(meta.OwnerTokenID)

	// 获取 tokenID
	receipt, err := DefaultPmContractHelper().Transfer(tokenID, from, to)
	if err != nil {
		return common.Hash{}, errors.WithMessage(err, "Failed to transfer owner")
	}

	return receipt.TransactionHash, nil
}

func (o *OwnerOperator) CheckIsContentOwner(account common.Address, name string) (bool, error) {
	meta, err := GetContentMetadata(name)
	if err != nil {
		return false, errors.WithMessage(err, "Failed to get content metadata")
	}

	owner, err := DefaultPmContractHelper().OwnerOf(bigutils.MustParseBigInt(meta.OwnerTokenID))
	if err != nil {
		return false, errors.WithMessage(err, "Failed to check if owner")
	}

	return owner == account, nil
}
