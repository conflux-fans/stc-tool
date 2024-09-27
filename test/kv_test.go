package test

import (
	"context"
	"fmt"
	"math"
	"os"
	"testing"
	"time"

	zg_common "github.com/0glabs/0g-storage-client/common"
	"github.com/0glabs/0g-storage-client/common/blockchain"
	ccore "github.com/0glabs/0g-storage-client/core"
	"github.com/0glabs/0g-storage-client/kv"
	"github.com/0glabs/0g-storage-client/node"
	"github.com/0glabs/0g-storage-client/transfer"
	"github.com/conflux-fans/storage-cli/constants"
	"github.com/conflux-fans/storage-cli/logger"
	"github.com/ethereum/go-ethereum/common"
	providers "github.com/openweb3/go-rpc-provider/provider_wrapper"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestKvGet(t *testing.T) {
	// new kv client
	zgsClient := node.MustNewZgsClient("http://127.0.0.1:15000", providers.Option{
		Logger: os.Stdout,
	})

	blockchainClient := blockchain.MustNewWeb3("http://127.0.0.1:14000", "7c5da44cf462b81e0b61a582f8c9b23ca78fc23e7104138f4e4329a9b2076e23", providers.Option{
		Logger: os.Stdout,
	})
	blockchain.CustomGasLimit = 1000000

	kvBatcher := kv.NewBatcher(math.MaxUint64, []*node.ZgsClient{zgsClient}, blockchainClient, zg_common.LogOption{
		Logger: logrus.New(),
	})
	kvClientForGet := kv.NewClient(node.MustNewKvClient("http://127.0.0.1:16000")) //

	// upload a key with streamID
	streamID := common.HexToHash("0x000000000000000000000000000000000000000000000000000000000000f2bd")
	batcher := kvBatcher
	batcher.Set(streamID,
		[]byte("TESTKEY0"),
		[]byte{69, 70, 71, 72, 73},
	)
	_, err := batcher.Exec(context.Background())
	assert.NoError(t, err)

	// get
	iter := kvClientForGet.NewIterator(streamID)
	fmt.Println("begin to end:")
	iter.SeekToFirst(context.Background())
	for iter.Valid() {
		pair := iter.KeyValue()
		fmt.Printf("%v: %v\n", string(pair.Key), string(pair.Data))
		iter.Next(context.Background())
	}

	value, err := kvClientForGet.GetValue(context.Background(), streamID, []byte("TESTKEY0"))
	assert.NoError(t, err)
	assert.Equal(t, []byte{69, 70, 71, 72, 73}, value.Data)
}

func TestZgUploadBytes(t *testing.T) {
	w3client := blockchain.MustNewWeb3("http://127.0.0.1:14000", "7c5da44cf462b81e0b61a582f8c9b23ca78fc23e7104138f4e4329a9b2076e23", providers.Option{
		Logger: os.Stdout,
	})
	zgNodeClients := node.MustNewZgsClients([]string{"http://127.0.0.1:15000", "http://127.0.0.1:15001"}, providers.Option{
		Logger: os.Stdout,
	})

	uploader, err := transfer.NewUploader(context.Background(), w3client, zgNodeClients, zg_common.LogOption{
		Logger: logrus.New(),
	})
	assert.NoError(t, err)
	dataInMemory, err := ccore.NewDataInMemory([]byte("hello world" + time.Now().String()))
	assert.NoError(t, err)

	_, err = uploader.Upload(context.Background(), dataInMemory)
	assert.NoError(t, err)
}

func TestKvBatchUploadBytes(t *testing.T) {
	kvStreamId := common.HexToHash("000000000000000000000000000000000000000000000000000000000000f2bd")

	w3client := blockchain.MustNewWeb3("http://127.0.0.1:14000", "7c5da44cf462b81e0b61a582f8c9b23ca78fc23e7104138f4e4329a9b2076e23", providers.Option{
		Logger: os.Stdout,
	})
	zgNodeClients := node.MustNewZgsClients([]string{"http://127.0.0.1:15000"}, providers.Option{
		Logger: os.Stdout,
	})
	kvBatcher := kv.NewBatcher(math.MaxUint64, zgNodeClients, w3client, zg_common.LogOption{
		Logger: logrus.New(),
	})

	now := time.Now()
	dataInMemory, err := ccore.NewDataInMemory([]byte("hello world"))
	assert.NoError(t, err)

	iterator := dataInMemory.Iterate(0, int64(constants.CHUNK_SIZE), false)
	i := 0
	for {
		exist, err := iterator.Next()
		assert.NoError(t, err)
		if !exist {
			break
		}
		chunk := iterator.Current()
		key := []byte(fmt.Sprintf("%v:%d", now, i))
		kvBatcher.Set(kvStreamId, key, chunk)
		logger.Get().WithField("key", string(key)).Info("Set line kv")
		i++
	}
	_, err = kvBatcher.Exec(context.Background())
	assert.NoError(t, err)

	kvClientForGet := kv.NewClient(node.MustNewKvClient("http://127.0.0.1:16000")) //

	iter := kvClientForGet.NewIterator(kvStreamId)
	fmt.Println("begin to end:")
	iter.SeekToFirst(context.Background())
	for iter.Valid() {
		pair := iter.KeyValue()
		fmt.Printf("%v: %v\n", string(pair.Key), string(pair.Data))
		iter.Next(context.Background())
	}

	value, err := kvClientForGet.GetValue(context.Background(), kvStreamId, []byte(fmt.Sprintf("%v:%d", now, 0)))
	assert.NoError(t, err)
	assert.Equal(t, "hello world", value.Data)
}

func TestZgDownload(t *testing.T) {
	root := common.HexToHash("0x276b14f314e7d3583c6718c75f8fc2e1d89b0f13446bf1ee5a02ab8457325343")
	// w3client := blockchain.MustNewWeb3("http://127.0.0.1:14000", "7c5da44cf462b81e0b61a582f8c9b23ca78fc23e7104138f4e4329a9b2076e23", providers.Option{
	// 	Logger: os.Stdout,
	// })
	zgNodeClients := node.MustNewZgsClients([]string{"http://127.0.0.1:15000", "http://127.0.0.1:15001"}, providers.Option{
		Logger: os.Stdout,
	})

	downloader, err := transfer.NewDownloader(zgNodeClients, zg_common.LogOption{
		Logger: logrus.New(),
	})
	assert.NoError(t, err)

	// _, err = uploader.Upload(context.Background(), dataInMemory)
	// assert.NoError(t, err)
	err = downloader.Download(context.Background(), root.Hex(), "../tmp/test_download.txt", false)
	assert.NoError(t, err)
}
