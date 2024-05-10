package test

import (
	"fmt"
	"testing"

	"github.com/0glabs/0g-storage-client/common/blockchain"
	"github.com/0glabs/0g-storage-client/contract"
	"github.com/0glabs/0g-storage-client/kv"
	"github.com/0glabs/0g-storage-client/node"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestKvUpdate(t *testing.T) {
	// new kv client
	zgsClient, err := node.NewClient("http://127.0.0.1:11100")
	if err != nil {
		fmt.Println(err)
		return
	}
	blockchainClient := blockchain.MustNewWeb3("https://evmtestnet.confluxrpc.com/6XWH2kDUX4wcKVN1VThMpjhhwerkTMZR8GYjk3S8Ti6GhM8qw7TJXDuT4sJWsM8MNmz2oxLsWAbjDUELaeAG4QA9Y", "7c5da44cf462b81e0b61a582f8c9b23ca78fc23e7104138f4e4329a9b2076e23")
	defer blockchainClient.Close()
	blockchain.CustomGasLimit = 1000000
	zgs, err := contract.NewFlowContract(common.HexToAddress("0x7eAce3C279B9fb4CCCE449d718e471fc2e330ea5"), blockchainClient)
	if err != nil {
		fmt.Println(err)
		return
	}
	kvClient := kv.NewClient(zgsClient, zgs)

	// upload a key with streamID
	streamID := common.HexToHash("0x000000000000000000000000000000000000000000000000000000000000f2bd")
	batcher := kvClient.Batcher()
	batcher.Set(streamID,
		[]byte("TESTKEY0"),
		[]byte{69, 70, 71, 72, 73},
	)
	err = batcher.Exec()
	assert.NoError(t, err)

	// get
	iter := kvClient.NewIterator(streamID)
	fmt.Println("begin to end:")
	iter.SeekToFirst()
	for iter.Valid() {
		pair := iter.KeyValue()
		fmt.Printf("%v: %v\n", string(pair.Key), string(pair.Data))
		iter.Next()
	}

	// upload another key
	batcher.Set(streamID,
		[]byte("TESTKEY1"),
		[]byte{74, 75, 76, 77, 78},
	)
	// get

}
