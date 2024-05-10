package core

import (
	"github.com/0glabs/0g-storage-client/node"
	"github.com/ethereum/go-ethereum/common"
)

func GetFileInfo(root common.Hash) (*node.FileInfo, error) {
	return nodeClients[0].ZeroGStorage().GetFileInfo(root)
}
