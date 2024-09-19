package test

import (
	"context"
	"fmt"
	"math"
	"os"
	"testing"

	zg_common "github.com/0glabs/0g-storage-client/common"
	"github.com/0glabs/0g-storage-client/common/blockchain"
	"github.com/0glabs/0g-storage-client/kv"
	"github.com/0glabs/0g-storage-client/node"
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
