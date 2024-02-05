package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/zero-gravity-labs/zerog-storage-client/node"
)

func GetFileInfo(root common.Hash) (*node.FileInfo, error) {
	return nodeClients[0].ZeroGStorage().GetFileInfo(root)
}
