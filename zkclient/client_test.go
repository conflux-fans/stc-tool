package zkclient

import (
	"fmt"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	providers "github.com/openweb3/go-rpc-provider/provider_wrapper"
	"github.com/openweb3/web3go"
	"github.com/stretchr/testify/assert"
)

func TestGenProof(t *testing.T) {
	c := MustNewClientWithOption("http://127.0.0.1:3030", web3go.ClientOption{
		Option: providers.Option{
			Logger: os.Stdout,
		},
	})

	vc := VC{
		Name:      "Alice",
		Age:       25,
		BirthDate: "20000101",
		EduLevel:  4,
		SerialNo:  "1234567890",
	}

	p, err := c.GetProof(&ProveInput{
		Data:               vc,
		BirthdateThreshold: "20110101",
		MerkleProof:        []common.Hash{common.Hash{}, common.Hash{}, common.Hash{}},
		PathIndex:          0,
	})

	assert.NoError(t, err)

	fmt.Printf("p %v\n", p)

	c.Verify(p, VerifyInput{})
}
