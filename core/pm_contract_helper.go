package core

import (
	"fmt"
	"math/big"
	"time"

	"github.com/conflux-fans/storage-cli/config"
	"github.com/conflux-fans/storage-cli/contracts"
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

	client := web3go.MustNewClientWithOption(config.Get().BlockChain.URL, web3go.ClientOption{
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

func (o *PmContractHelper) getPmContract(from common.Address) (*contracts.PermissionManager, bind.SignerFn, error) {
	client, ok := w3Clients[from]
	if !ok {
		return nil, nil, fmt.Errorf("w3client of account %s not found", from.Hex())
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
