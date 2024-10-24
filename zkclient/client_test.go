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
		common.HexToHash("0xcca4353c87fe7c16438b81f6931e04112daa4c6c880e0df49a8950f27b4bcc23"),
		common.HexToHash("0xf73e6947d7d1628b9976a6e40d7b278a8a16405e96324a68df45b12a51b7cfde"),
		common.HexToHash("0x697dbb83b0e50fa3961362f3ca30f855bad9357b1e6b77ffb1851a891882a59d"),
		common.HexToHash("0xde5747106ac1194a1fa9071dbd6cf19dc2bc7964497ef0afec7e4bdbcf08c47e"),
		common.HexToHash("0x09c7082879180d28c789c05fafe7030871c76cedbe82c948b165d6a1d66ac15b"),
		common.HexToHash("0xaa7a02bcf29fba687f84123c808b5b48834ff5395abe98e622fadc14e4180c95"),
		common.HexToHash("0x7608fd46b710b589e0f2ee5a13cd9c41d432858a30d524f84c6d5db37f66273a"),
		common.HexToHash("0xa5d9a2f7f3573ac9a1366bc484688b4daf934b87ea9b3bf2e703da8fd9f09708"),
		common.HexToHash("0x6c1779477f4c3fca26b4607398859a43b90a286ce8062500744bd4949981757f"),
		common.HexToHash("0x45c22df3d952c33d5edce122eed85e5cda3fd61939e7ad7b3e03b6927bb598ea"),
		common.HexToHash("0xf56da5d89b3daa7083544e415e4956a43838e4c6115265aeb21fd63543e8b2cb"),
		common.HexToHash("0x62d78399b954d51cb9728601738ad13ddc43b2300064660716bb661d2f4d686f"),
		common.HexToHash("0xd5ed03fe35fa079b4b6ae3ffd5dc56af3e85f54a44415b1cf00e94de41bd8107"),
		common.HexToHash("0x1d1a3a74062fd94078617e33eb901eaf16a830f67c387d8eed342db2ac5e2cc5"),
		common.HexToHash("0x19b3b3886526917eae8650223d0be20a0301be960eb339696e673ad8a804440f"),
		common.HexToHash("0xc5d8c52ca4f9330fcd2eace964f1c04a7baa8053b3795e8c738c6e3901f166ed"),
		common.HexToHash("0x0247df4aaebfe15af056de45b0b5a286b025afae9f70f4eb6c2a0aebe8cc06c3"),
	}, 103429, []ExtensionSignal{{Date: &birthDateThreshold}}))

	assert.NoError(t, err)

	fmt.Printf("get proof: %v\n", p)
}

func TestVerifyProof(t *testing.T) {
	c := MustNewClientWithOption("http://127.0.0.1:3030", web3go.ClientOption{
		Option: providers.Option{
			Logger:         os.Stdout,
			RequestTimeout: time.Hour,
		},
	})

	birthDateThreshold := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	proof := "0679c865c006364eef6f7241dd7280b763bd92faef337b224a8d784cf43f058ca80cc0b772f7609a8a91bf3c21f72ccd133ee11c5212caaa8d00fddac0262c099602c78b75db2d671a6e4ed13a0abad8cfd7996eceb56e24d052de0cacf7e1801e5d7eb95bb5ea63aea0706b2b5eff9f70925a54482bb38c922df50bdae95498"
	verify, err := c.Verify(proof, VerifyInput{
		Extensions: []ExtensionSignal{{Date: &birthDateThreshold}},
		Root:       common.HexToHash("0x717aaedd8d94331647ae8c81acf986b2ad5413237e3246ce6dfe009a2d809627"),
	})
	assert.NoError(t, err)
	assert.Equal(t, true, verify)
	fmt.Printf("verify result: %v\n", verify)
}
