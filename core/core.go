package core

import (
	"fmt"
	"os"

	"github.com/conflux-fans/storage-cli/config"
	"github.com/conflux-fans/storage-cli/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	providers "github.com/openweb3/go-rpc-provider/provider_wrapper"
	"github.com/openweb3/web3go/signers"
	"github.com/sirupsen/logrus"
	"github.com/zero-gravity-labs/zerog-storage-client/common/blockchain"
	"github.com/zero-gravity-labs/zerog-storage-client/contract"
	"github.com/zero-gravity-labs/zerog-storage-client/kv"
	"github.com/zero-gravity-labs/zerog-storage-client/node"
)

var (
	// w3client              *web3go.Client
	nodeClients         []*node.Client
	kvClientForIterator *kv.Client
	kvClientsForPut     map[common.Address]*kv.Client
	accounts            []common.Address
	adminKvClientForPut *kv.Client
	defaultFlow         *contract.FlowContract
	templates           *contracts.Templates
	defaultAccount      common.Address
	signerFn            bind.SignerFn
	// signerManager         *signers.SignerManager
)

var (
	STREAM_FILE = common.HexToHash("000000000000000000000000000000000000000000000000000000000000f2bd")
)

func Init() {
	cfg := config.Get()
	// logrus.WithField("config", cfg).Info("Get config")
	nodeClients = node.MustNewClients(cfg.StorageNodes)

	providerOpt := providers.Option{}
	if cfg.Log == config.DEBUG {
		providerOpt.Logger = os.Stdout
	}
	kvClientForIterator = kv.NewClient(node.MustNewClient(cfg.KvNode, providerOpt), defaultFlow)

	genKvClientsForPut()
	initTempalteContract()
	grantAllAccountWriter()
}

func genKvClientsForPut() {
	kvClientsForPut = make(map[common.Address]*kv.Client)
	cfg := config.Get()

	for i, pk := range cfg.PrivateKeys {
		w3client := blockchain.MustNewWeb3(cfg.BlockChain.URL, pk)
		if cfg.Log == config.DEBUG {
			w3client.SetProvider(providers.NewLoggerProvider(w3client.Provider(), os.Stdout))
		}

		flowAddr := common.HexToAddress(cfg.BlockChain.FlowContract)
		flow, err := contract.NewFlowContract(flowAddr, w3client)
		if err != nil {
			logrus.WithError(err).Fatal("Failed to create flow contract")
			os.Exit(1)
		}
		kvClient := kv.NewClient(nodeClients[0], flow)
		account := signers.MustNewPrivateKeySignerByString(pk).Address()
		accounts = append(accounts, account)
		kvClientsForPut[account] = kvClient
		if i == 0 {
			adminKvClientForPut = kvClient
			defaultFlow = flow
			defaultAccount = account
		}
	}
}

func initTempalteContract() {
	cfg := config.Get()
	w3client := blockchain.MustNewWeb3(cfg.BlockChain.URL, cfg.PrivateKeys[0])
	templateAddr := common.HexToAddress(cfg.BlockChain.TemplateContract)
	backend, _signerFn := w3client.ToClientForContract()
	signerFn = _signerFn

	var err error
	templates, err = contracts.NewTemplates(templateAddr, backend)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create templates contract")
		os.Exit(1)
	}
}

func grantAllAccountWriter() {
	if err := GrantStreamWriter(accounts...); err != nil {
		panic(fmt.Sprintf("Failed to grant account %v strem writer", accounts))
	}
}
