package zkclient

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	providers "github.com/openweb3/go-rpc-provider/provider_wrapper"
	"github.com/openweb3/web3go"
	"github.com/stretchr/testify/assert"
)

func TestGenProof(t *testing.T) {
	c := MustNewClientWithOption("http://127.0.0.1:3030", web3go.ClientOption{
		Option: providers.Option{
			Logger:         os.Stdout,
			RequestTimeout: time.Hour,
		},
	})

	vc := VC{
		Name:      "Alice",
		Age:       25,
		BirthDate: "20000101",
		EduLevel:  4,
		SerialNo:  "1234567890",
	}

	birthDateThreshold := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	p, err := c.GetProof(NewProveInput("verysecretkey123", "uniqueiv12345678", vc, []common.Hash{
		common.HexToHash("0x9c8e52ca76f7555771149cbcd84372549c8e52ca76f7555771149cbcd8437254"),
		common.HexToHash("0x2ddf90a82346f547b58f7804e3ea9a112ddf90a82346f547b58f7804e3ea9a11"),
		common.HexToHash("0x0ecb5efcadfe3664b95fbfd68fe560e60ecb5efcadfe3664b95fbfd68fe560e6"),
	}, 0, []ExtensionSignal{{Date: &birthDateThreshold}}))

	assert.NoError(t, err)

	fmt.Printf("get proof: %v\n", p)

	verify, err := c.Verify(p, VerifyInput{
		Extensions: []ExtensionSignal{{Date: &birthDateThreshold}},
		Root:       common.HexToHash("0x9c8e52ca76f7555771149cbcd84372549c8e52ca76f7555771149cbcd8437254"),
	})
	assert.NoError(t, err)
	fmt.Printf("verify result: %v\n", verify)
}
