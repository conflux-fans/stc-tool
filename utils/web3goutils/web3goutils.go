package web3goutils

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/openweb3/web3go"
	"github.com/openweb3/web3go/types"
)

func WaitTransactionReceipt(timeout time.Duration, client *web3go.Client, txHash common.Hash, pollInterval time.Duration) (*types.Receipt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			receipt, err := client.Eth.TransactionReceipt(txHash)
			if err != nil {
				return nil, err
			}
			if receipt != nil {
				return receipt, nil
			}
			time.Sleep(pollInterval)
		}
	}
}
