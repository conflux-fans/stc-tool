package core

import (
	"errors"
	"fmt"
	"math"
	"os"
	"sync"
	"time"

	"github.com/0glabs/0g-storage-client/common/blockchain"
	"github.com/0glabs/0g-storage-client/contract"
	"github.com/0glabs/0g-storage-client/kv"
	"github.com/0glabs/0g-storage-client/node"
	"github.com/conflux-fans/storage-cli/config"
	"github.com/conflux-fans/storage-cli/contracts"
	"github.com/conflux-fans/storage-cli/logger"
	"github.com/conflux-fans/storage-cli/zkclient"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	providers "github.com/openweb3/go-rpc-provider/provider_wrapper"
	"github.com/openweb3/web3go"
	"github.com/openweb3/web3go/signers"
)

var (
	// w3client              *web3go.Client
	zgNodeClients       []*node.ZgsClient
	kvClientForIterator *kv.Client
	kvBatcherForPut     map[common.Address]*kv.Batcher
	w3Clients           map[common.Address]*web3go.Client
	zkClient            *zkclient.Client
	accounts            []common.Address
	adminBatcher        *kv.Batcher
	adminW3Client       *web3go.Client
	defaultFlow         *contract.FlowContract
	templates           *contracts.Templates
	defaultAccount      common.Address
	signerFn            bind.SignerFn
	// signerManager         *signers.SignerManager
	grantWriterOnce sync.Once
)

var (
	STREAM_FILE = common.HexToHash("000000000000000000000000000000000000000000000000000000000000f2bd")
)

func Init() {
	cfg := config.Get()
	// logger.Get().WithField("config", cfg).Info("Get config")
	zgNodeClients = node.MustNewZgsClients(cfg.StorageNodes)
	zkClient = zkclient.MustNewClientWithOption(cfg.ZkNode, web3go.ClientOption{
		Option: providers.Option{
			Logger:         os.Stdout,
			RequestTimeout: time.Minute,
		},
	})

	providerOpt := providers.Option{}
	if cfg.Log == config.DEBUG {
		providerOpt.Logger = os.Stdout
	}
	kvClientForIterator = kv.NewClient(node.MustNewKvClient(cfg.KvNode, providerOpt))

	genKvClientsForPut()
	initTempalteContract()
	// GrantAllAccountWriter()
}

func genKvClientsForPut() {
	kvBatcherForPut = make(map[common.Address]*kv.Batcher)
	cfg := config.Get()

	for i, pk := range cfg.PrivateKeys {
		w3client := blockchain.MustNewWeb3(cfg.BlockChain.URL, pk)
		if cfg.Log == config.DEBUG {
			w3client.SetProvider(providers.NewLoggerProvider(w3client.Provider(), os.Stdout))
		}

		flowAddr := common.HexToAddress(cfg.BlockChain.FlowContract)
		flow, err := contract.NewFlowContract(flowAddr, w3client)
		if err != nil {
			logger.Get().WithError(err).Fatal("Failed to create flow contract")
			os.Exit(1)
		}
		// kvClient := kv.NewClient(zgNodeClients[0])
		kvBatcher := kv.NewBatcher(math.MaxInt64, zgNodeClients, w3client)
		account := signers.MustNewPrivateKeySignerByString(pk).Address()
		accounts = append(accounts, account)
		kvBatcherForPut[account] = kvBatcher
		if i == 0 {
			adminBatcher = kvBatcher
			adminW3Client = w3client
			defaultFlow = flow
			defaultAccount = account
		}
	}
}

func getKvClientBatcher(account common.Address) (*kv.Batcher, error) {
	if kvBatcherForPut[account] == nil {
		return nil, errors.New("no kv client for account")
	}
	return kvBatcherForPut[account], nil
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
		logger.Get().WithError(err).Fatal("Failed to create templates contract")
		os.Exit(1)
	}
}

func GrantAllAccountStreamWriter() {
	grantWriterOnce.Do(func() {
		if err := GrantStreamWriter(accounts...); err != nil {
			panic(fmt.Sprintf("Failed to grant account %v strem writer", accounts))
		}
	})
	// logger.Get().Info("Granted all")
}
