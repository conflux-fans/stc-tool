package core

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/conflux-fans/storage-cli/config"
	"github.com/conflux-fans/storage-cli/logger"
	"github.com/conflux-fans/storage-cli/pkg/utils/bigutils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/assert"
)

func TestFilterTransfer(t *testing.T) {
	config.SetConfigFile("/Users/dayong/myspace/mywork/storage-tool/config.yaml")
	config.Init()
	Init()

	latest, err := adminW3Client.Eth.BlockNumber()
	assert.NoError(t, err)
	fmt.Printf("latest: %d\n", latest)

	start := uint64(latest.Int64()) - 200000
	end := uint64(latest.Int64())
	transfers, err := DefaultPmContractHelper().FilterTransfer(&bind.FilterOpts{Start: start, End: &end}, nil, nil, []*big.Int{bigutils.MustParseBigInt("1")})
	assert.NoError(t, err)
	j, _ := json.Marshal(transfers)
	logger.Get().WithField("transfers", string(j)).Info("transfers")
}
