package core

import (
	"math/big"

	"github.com/conflux-fans/storage-cli/config"
	"github.com/conflux-fans/storage-cli/logger"
	"github.com/conflux-fans/storage-cli/pkg/utils/bigutils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

var ownerOperator OwnerOperator

type OwnerOperator struct {
}

func DefaultOwnerOperator() *OwnerOperator {
	return &ownerOperator
}

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

func (o *OwnerOperator) GetOwnerHistory(name string) ([]common.Address, error) {
	meta, err := GetContentMetadata(name)
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to get content metadata")
	}
	logger.Get().WithField("token_id", meta.OwnerTokenID).Info("get content owner token id")

	opts, err := getFilterOpt()
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to get filter options")
	}

	logs, err := DefaultPmContractHelper().FilterTransfer(opts, nil, nil, []*big.Int{bigutils.MustParseBigInt(meta.OwnerTokenID)})
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to get owner history")
	}

	var result []common.Address
	for i, log := range logs {
		if i > 0 && logs[i-1].To != log.From {
			return nil, errors.New("Owner history is not continuous")
		}
		result = append(result, log.From)
	}

	return result, nil
}

func getFilterOpt() (*bind.FilterOpts, error) {
	latestBlock, err := adminW3Client.Eth.BlockNumber()
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to get latest block number")
	}
	start := uint64(config.Get().BlockChain.StartBlockNum)
	end := uint64(latestBlock.Int64())
	opts := &bind.FilterOpts{
		Start: start,
		End:   &end,
	}
	return opts, nil
}
