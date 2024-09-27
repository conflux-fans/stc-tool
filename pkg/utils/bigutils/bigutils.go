package bigutils

import (
	"math/big"
)

func MustParseBigInt(tokenID string) *big.Int {
	tokenIDInt, ok := new(big.Int).SetString(tokenID, 10)
	if !ok {
		panic("invalid big.Int")
	}
	return tokenIDInt
}
