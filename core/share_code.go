package core

import (
	"encoding/base64"

	"github.com/ethereum/go-ethereum/common"
)

type ShareCodeHelper struct {
}

func NewShareCodeHelper() *ShareCodeHelper {
	return &ShareCodeHelper{}
}

func (s *ShareCodeHelper) GetShareCode(root common.Hash) string {
	return base64.StdEncoding.EncodeToString(root.Bytes())
}

func (s *ShareCodeHelper) GetRootFromShareCode(shareCode string) (common.Hash, error) {
	decoded, err := base64.StdEncoding.DecodeString(shareCode)
	if err != nil {
		return common.Hash{}, err
	}
	return common.BytesToHash(decoded), nil
}
