package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/0glabs/0g-storage-client/common/blockchain"
	"github.com/0glabs/0g-storage-client/contract"
	"github.com/0glabs/0g-storage-client/kv"
	"github.com/0glabs/0g-storage-client/node"
	"github.com/conflux-fans/storage-cli/core"
	"github.com/ethereum/go-ethereum/common"
	providers "github.com/openweb3/go-rpc-provider/provider_wrapper"
	"github.com/stretchr/testify/assert"
)

func TestKvUpdate(t *testing.T) {
	// new kv client
	zgsClient := node.MustNewClient("http://127.0.0.1:11100")
	blockchainClient := blockchain.MustNewWeb3("https://evmtestnet.confluxrpc.com/6XWH2kDUX4wcKVN1VThMpjhhwerkTMZR8GYjk3S8Ti6GhM8qw7TJXDuT4sJWsM8MNmz2oxLsWAbjDUELaeAG4QA9Y", "7c5da44cf462b81e0b61a582f8c9b23ca78fc23e7104138f4e4329a9b2076e23")
	defer blockchainClient.Close()
	blockchain.CustomGasLimit = 1000000
	zgs, err := contract.NewFlowContract(common.HexToAddress("0x585baAC11508a326Bb2aaf5a19Ad08AD8f98aD4f"), blockchainClient)
	if err != nil {
		fmt.Println(err)
		return
	}
	kvClientForPut := kv.NewClient(zgsClient, zgs)
	kvClientForGet := kv.NewClient(node.MustNewClient("http://127.0.0.1:6789"), zgs)

	// upload a key with streamID
	streamID := common.HexToHash("0x000000000000000000000000000000000000000000000000000000000000f2bd")
	batcher := kvClientForPut.Batcher()
	batcher.Set(streamID,
		[]byte("TESTKEY0"),
		[]byte{69, 70, 71, 72, 73},
	)
	err = batcher.Exec()
	assert.NoError(t, err)

	// get
	iter := kvClientForGet.NewIterator(streamID)
	fmt.Println("begin to end:")
	iter.SeekToFirst()
	for iter.Valid() {
		pair := iter.KeyValue()
		fmt.Printf("%v: %v\n", string(pair.Key), string(pair.Data))
		iter.Next()
	}

	value, err := kvClientForPut.GetValue(streamID, []byte("TESTKEY0"))
	assert.NoError(t, err)
	assert.Equal(t, []byte{69, 70, 71, 72, 73}, value)
}

func TestKvGet(t *testing.T) {
	// new kv client
	zgsClient, err := node.NewClient("http://127.0.0.1:11100", providers.Option{
		Logger: os.Stdout,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	blockchainClient := blockchain.MustNewWeb3("https://evmtestnet.confluxrpc.com/6XWH2kDUX4wcKVN1VThMpjhhwerkTMZR8GYjk3S8Ti6GhM8qw7TJXDuT4sJWsM8MNmz2oxLsWAbjDUELaeAG4QA9Y", "7c5da44cf462b81e0b61a582f8c9b23ca78fc23e7104138f4e4329a9b2076e23")
	defer blockchainClient.Close()
	blockchain.CustomGasLimit = 1000000
	zgs, err := contract.NewFlowContract(common.HexToAddress("0x585baAC11508a326Bb2aaf5a19Ad08AD8f98aD4f"), blockchainClient)
	if err != nil {
		fmt.Println(err)
		return
	}
	kvClientForSet := kv.NewClient(zgsClient, zgs)

	// kvClientForGet := kv.NewClient(zgsClient, zgs)

	// // get
	// iter := kvClient.NewIterator(core.STREAM_FILE)
	// fmt.Println("begin to end:")
	// iter.SeekToFirst()
	// for iter.Valid() {
	// 	pair := iter.KeyValue()
	// 	fmt.Printf("%v: %v\n", string(pair.Key), string(pair.Data))
	// 	iter.Next()
	// }

	batcher := kvClientForSet.Batcher()
	batcher.Set(core.STREAM_FILE,
		[]byte("TESTKEY0"),
		[]byte{69, 70, 71, 72, 73},
	)
	err = batcher.Exec()
	assert.NoError(t, err)

	value, err := kvClientForSet.GetValue(core.STREAM_FILE, []byte("TESTKEY0"))
	assert.NoError(t, err)
	assert.Equal(t, []byte{69, 70, 71, 72, 73}, value)
}
