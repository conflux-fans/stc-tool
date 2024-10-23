package core

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/0glabs/0g-storage-client/node"
	"github.com/conflux-fans/storage-cli/config"
	"github.com/conflux-fans/storage-cli/zkclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestConvertFlowProofToForZk(t *testing.T) {
	nodeProof := node.FlowProof{
		Lemma: []common.Hash{
			common.HexToHash("0xcca4353c87fe7c16438b81f6931e04112daa4c6c880e0df49a8950f27b4bcc23"),
			common.HexToHash("0xcca4353c87fe7c16438b81f6931e04112daa4c6c880e0df49a8950f27b4bcc23"),
			common.HexToHash("0xf73e6947d7d1628b9976a6e40d7b278a8a16405e96324a68df45b12a51b7cfde"),
			common.HexToHash("0x697dbb83b0e50fa3961362f3ca30f855bad9357b1e6b77ffb1851a891882a59d"),
			common.HexToHash("0xde5747106ac1194a1fa9071dbd6cf19dc2bc7964497ef0afec7e4bdbcf08c47e"),
		},
		Path: []bool{
			true, false, true,
		},
	}
	zkProof := convertFlowProofToForZk(&nodeProof)
	assert.Equal(t, zkProof.Lemma, nodeProof.Lemma[1:4])
	assert.Equal(t, zkProof.Path, uint64(2))
	fmt.Printf("zkProof: %v\n", zkProof)
}

func TestZk(t *testing.T) {
	os.Chdir("/Users/dayong/myspace/mywork/storage-tool")
	config.Init()
	Init()

	vc := `{"name": "Alice", "age": 25, "birth_date": "20000101", "edu_level": 4, "serial_no": "1234567890"}`
	var _vc zkclient.VC
	err := json.Unmarshal([]byte(vc), &_vc)
	assert.NoError(t, err)

	ZkUploadInput := ZkUploadInput{
		Vc:                 &_vc,
		BirthdateThreshold: "20000101",
	}
	key, iv := "verysecretkey123", "uniqueiv12345678"

	t.Run("integration", func(_t *testing.T) {
		// upload
		zkUploadOutput, err := NewZk().UploadVc(&ZkUploadInput, key, iv)
		assert.NoError(_t, err)
		fmt.Printf("zk upload output: %v\n", zkUploadOutput)
		assert.Equal(_t, zkUploadOutput.VcDataRoot.Hex(), "0xcca4353c87fe7c16438b81f6931e04112daa4c6c880e0df49a8950f27b4bcc23")

		// proof
		flowProofForZk := FlowProofForZk{
			VcDataRoot: zkUploadOutput.VcDataRoot,
			FlowRoot:   zkUploadOutput.FlowRoot,
			Lemma:      zkUploadOutput.Lemma,
			Path:       zkUploadOutput.Path,
		}
		proveInput := &ZkProofInput{
			ZkUploadInput:  ZkUploadInput,
			FlowProofForZk: flowProofForZk,
			Key:            key,
			IV:             iv,
		}
		proveOutput, err := NewZk().ZkProof(proveInput)
		assert.NoError(_t, err)
		fmt.Printf("prove output: %v\n", proveOutput)

		// verify
		// Note: 无法直接测试正确性，因为在生成 proof 时有随机数参与，所以会导致每次结果不一样；只能通过 verify 测试
		isSucess, err := NewZk().ZkVerify(proveOutput.Proof, "20000101", flowProofForZk.FlowRoot.Hex())
		assert.NoError(t, err)
		assert.True(t, isSucess)
		fmt.Println(isSucess)
	})

	t.Run("upload vc", func(_t *testing.T) {
		zkUploadOutput, err := NewZk().UploadVc(&ZkUploadInput, key, iv)
		assert.NoError(_t, err)
		fmt.Printf("zk upload output: %v\n", zkUploadOutput)
	})

	t.Run("zk proof", func(_t *testing.T) {
		flowProofForZk := FlowProofForZk{
			VcDataRoot: common.HexToHash("0xcca4353c87fe7c16438b81f6931e04112daa4c6c880e0df49a8950f27b4bcc23"),
			FlowRoot:   common.HexToHash("0x717aaedd8d94331647ae8c81acf986b2ad5413237e3246ce6dfe009a2d809627"),
			Lemma: []common.Hash{
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
			},
			Path: 103429,
		}

		proveInput := &ZkProofInput{
			ZkUploadInput:  ZkUploadInput,
			FlowProofForZk: flowProofForZk,
			Key:            key,
			IV:             iv,
		}
		proveOutput, err := NewZk().ZkProof(proveInput)
		assert.NoError(_t, err)
		fmt.Printf("prove output: %v\n", proveOutput)

		// Note: 无法直接测试正确性，因为在生成 proof 时有随机数参与，所以会导致每次结果不一样；只能通过 verify 测试
		isSucess, err := NewZk().ZkVerify(proveOutput.Proof, "20000101", flowProofForZk.FlowRoot.Hex())
		assert.NoError(t, err)
		assert.True(t, isSucess)
		fmt.Println(isSucess)
	})

	t.Run("zk verify", func(_t *testing.T) {
		proof := "8a3cdd45b15d762268606e64c88823b013a38aa5d1c8df2d361d631c28e49e048653fee284998ed4f3ec4fdcaf4c129515c3eb8e8aa3dc3cd18381ccbe73bd16f490dcba0e1cb746308f9e622c34d565da4c06e3aa6c2b7048b4d03525773c91985759bceda0d0dfe7174e8697d0c1e5d0b0c96fc506c794e0a7c820bf9abd9f"
		flowRoot := common.HexToHash("0x717aaedd8d94331647ae8c81acf986b2ad5413237e3246ce6dfe009a2d809627")

		isSucess, err := NewZk().ZkVerify(proof, "20000101", flowRoot.Hex())
		assert.NoError(t, err)
		assert.True(t, isSucess)
		fmt.Println(isSucess)
	})
}
