package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func streamIdByName(name string) common.Hash {
	return crypto.Keccak256Hash([]byte(name))
}
