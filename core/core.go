package core

import (
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/openweb3/web3go"
	"github.com/sirupsen/logrus"
	"github.com/zero-gravity-labs/zerog-storage-client/common/blockchain"
	"github.com/zero-gravity-labs/zerog-storage-client/contract"
	"github.com/zero-gravity-labs/zerog-storage-client/node"
	"github.com/zero-gravity-labs/zerog-storage-tool/config"
)

var (
	w3client    *web3go.Client
	nodeClients []*node.Client
	flow        *contract.FlowContract
)

func init() {
	cfg := config.Get()

	w3client = blockchain.MustNewWeb3(cfg.BlockChain.URL, cfg.PrivateKeys[0])
	nodeClients = node.MustNewClients(cfg.StorageNodes)

	var err error
	contractAddr := common.HexToAddress(cfg.BlockChain.FlowContract)
	flow, err = contract.NewFlowContract(contractAddr, w3client)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create flow contract")
		os.Exit(1)
	}
}
