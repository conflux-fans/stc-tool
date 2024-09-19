package core

import (
	"context"

	"github.com/0glabs/0g-storage-client/node"
	"github.com/ethereum/go-ethereum/common"
)

func GetFileInfo(root common.Hash) (*node.FileInfo, error) {
	return zgNodeClients[0].GetFileInfo(context.Background(), root)
}
