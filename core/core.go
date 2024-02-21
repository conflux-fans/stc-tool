package core

import (
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/openweb3/web3go"
	"github.com/openweb3/web3go/signers"
	"github.com/sirupsen/logrus"
	"github.com/zero-gravity-labs/zerog-storage-client/common/blockchain"
	"github.com/zero-gravity-labs/zerog-storage-client/contract"
	"github.com/zero-gravity-labs/zerog-storage-client/node"
	"github.com/zero-gravity-labs/zerog-storage-tool/config"
	"github.com/zero-gravity-labs/zerog-storage-tool/contracts"
)

var (
	w3client       *web3go.Client
	nodeClients    []*node.Client
	flow           *contract.FlowContract
	templates      *contracts.Templates
	signerManager  *signers.SignerManager
	defaultAccount common.Address
	signerFn       bind.SignerFn
)

func init() {
	cfg := config.Get()

	w3client = blockchain.MustNewWeb3(cfg.BlockChain.URL, cfg.PrivateKeys[0])
	nodeClients = node.MustNewClients(cfg.StorageNodes)

	var err error
	signerManager, err = w3client.GetSignerManager()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to get signer manager")
		os.Exit(1)
	}

	defaultAccount = signerManager.List()[0].Address()

	flowAddr := common.HexToAddress(cfg.BlockChain.FlowContract)
	flow, err = contract.NewFlowContract(flowAddr, w3client)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create flow contract")
		os.Exit(1)
	}

	templateAddr := common.HexToAddress(cfg.BlockChain.TemplateContract)
	backend, _signerFn := w3client.ToClientForContract()
	signerFn = _signerFn

	templates, err = contracts.NewTemplates(templateAddr, backend)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create templates contract")
		os.Exit(1)
	}
}
