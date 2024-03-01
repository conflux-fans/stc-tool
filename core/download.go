package core

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/sirupsen/logrus"
	"github.com/zero-gravity-labs/zerog-storage-client/node"
	"github.com/zero-gravity-labs/zerog-storage-client/transfer"
)

func DownloadFile(root string, savePath string) {
	downloader := transfer.NewDownloader(nodeClients...)
	if err := downloader.Download(root, root, false); err != nil {
		logrus.WithField("root", root).WithError(err).Fatal("Failed to download file")
	}
}

func DownloadByKv(streamName string) {
	streamId := crypto.Keccak256Hash([]byte(streamName))
	iter := kvClientForPut.NewIterator(streamId)

	iter.SeekToFirst()
	var kvs []*node.KeyValue
	for iter.Valid() {
		pair := iter.KeyValue()
		fmt.Printf("%v: %v\n", string(pair.Key), string(pair.Data))
		kvs = append(kvs, pair)
		iter.Next()
	}

	for _, kv := range kvs {
		// download file
		logrus.WithField("file index", kv.Key).Info("Downloading file")
		os.MkdirAll(streamName, 0777)
		DownloadFile(hexutil.Encode(kv.Data), streamName+string(kv.Key))
	}
}
