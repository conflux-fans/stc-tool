package core

import (
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/conflux-fans/storage-cli/config"
	"github.com/conflux-fans/storage-cli/contracts"
	"github.com/conflux-fans/storage-cli/logger"
	"github.com/conflux-fans/storage-cli/utils/web3goutils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/openweb3/web3go"
	"github.com/openweb3/web3go/types"
	"github.com/pkg/errors"
)

type PmContractHelper struct {
	contractAddr      common.Address
	pmContractForRead *contracts.PermissionManager
}

func DefaultPmContractHelper() *PmContractHelper {
	return NewOwnerOperator(common.HexToAddress(config.Get().BlockChain.PmContract))
}

func NewOwnerOperator(contract common.Address) *PmContractHelper {
	val := &PmContractHelper{
		contractAddr: contract,
	}

	client := web3go.MustNewClientWithOption(config.Get().BlockChain.ConfuraUrl, web3go.ClientOption{
		Option: providerOpt,
	})
	clientForContract, _ := client.ToClientForContract()
	pmContract, err := contracts.NewPermissionManager(contract, clientForContract)
	if err != nil {
		panic(err)
	}
	val.pmContractForRead = pmContract
	return val
}

func (o *PmContractHelper) getPmContract(signerAddr common.Address) (*contracts.PermissionManager, bind.SignerFn, error) {
	client, ok := w3Clients[signerAddr]
	if !ok {
		return nil, nil, fmt.Errorf("w3client of account %s not found", signerAddr.Hex())
	}

	clientForContract, signerFn := client.ToClientForContract()
	pmContract, err := contracts.NewPermissionManager(o.contractAddr, clientForContract)
	if err != nil {
		panic(err)
	}
	return pmContract, signerFn, nil
}

func (o *PmContractHelper) Mint(to common.Address) (*types.Receipt, *big.Int, error) {
	contractOwner, err := o.pmContractForRead.Owner(nil)
	if err != nil {
		return nil, nil, err
	}

	pm, signerFn, err := o.getPmContract(contractOwner)
	if err != nil {
		return nil, nil, err
	}

	tx, err := pm.SafeMint(&bind.TransactOpts{From: contractOwner, Signer: signerFn}, to)
	if err != nil {
		return nil, nil, err
	}

	receipt, err := web3goutils.WaitTransactionReceipt(20*time.Second, adminW3Client, tx.Hash(), 2*time.Second)
	if err != nil {
		return nil, nil, err
	}

	transferLog, err := pm.ParseTransfer(*receipt.Logs[0].ToEthLog())
	if err != nil {
		return nil, nil, err
	}

	return receipt, transferLog.TokenId, nil
}

func (o *PmContractHelper) Transfer(tokenId *big.Int, from common.Address, to common.Address) (*types.Receipt, error) {
	tokenOwner, err := o.pmContractForRead.OwnerOf(nil, tokenId)
	if err != nil {
		return nil, err
	}

	if tokenOwner != from {
		return nil, errors.New("not the owner of token")
	}

	pm, signerFn, err := o.getPmContract(tokenOwner)
	if err != nil {
		return nil, err
	}

	tx, err := pm.SafeTransferFrom(&bind.TransactOpts{From: from, Signer: signerFn}, tokenOwner, to, tokenId)
	if err != nil {
		return nil, err
	}

	receipt, err := web3goutils.WaitTransactionReceipt(20*time.Second, adminW3Client, tx.Hash(), 2*time.Second)
	if err != nil {
		return nil, err
	}
	return receipt, nil
}

func (o *PmContractHelper) OwnerOf(tokenId *big.Int) (common.Address, error) {
	owner, err := o.pmContractForRead.OwnerOf(nil, tokenId)
	if err != nil {
		return common.Address{}, err
	}
	return owner, nil
}

func (o *PmContractHelper) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenIds []*big.Int) ([]*contracts.PermissionManagerTransfer, error) {
	batch := config.Get().BlockChain.GetLogsBatchSize
	var result []*contracts.PermissionManagerTransfer
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errChan = make(chan error, 1)

	batchCount := (*opts.End - opts.Start + uint64(batch) - 1) / uint64(batch)
	logger.Get().WithField("start block", opts.Start).WithField("end block", *opts.End).WithField("batchCount", batchCount).Info("ready to get logs concurrently")

	batchResults := make([][]*contracts.PermissionManagerTransfer, batchCount)

	for i := uint64(0); i < batchCount; i++ {
		wg.Add(1)
		go func(i uint64) {
			defer wg.Done()
			start := opts.Start + i*uint64(batch)
			end := start + uint64(batch) - 1
			if end > *opts.End {
				end = *opts.End
			}

			batchOpts := &bind.FilterOpts{
				Start: start,
				End:   &end,
			}

			batchIter, err := o.pmContractForRead.FilterTransfer(batchOpts, from, to, tokenIds)
			if err != nil {
				select {
				case errChan <- err:
				default:
				}
				return
			}
			defer batchIter.Close()

			var batchLogs []*contracts.PermissionManagerTransfer
			for batchIter.Next() {
				batchLogs = append(batchLogs, batchIter.Event)
			}

			mu.Lock()
			batchResults[i] = batchLogs
			mu.Unlock()
		}(i)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	if err := <-errChan; err != nil {
		panic(err)
	}

	for _, batchLogs := range batchResults {
		result = append(result, batchLogs...)
	}

	return result, nil
}
