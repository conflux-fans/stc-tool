package core

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/conflux-fans/storage-cli/config"
	"github.com/conflux-fans/storage-cli/zkclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestZkUpload(t *testing.T) {
	os.Chdir("/Users/dayong/myspace/mywork/storage-tool")
	config.Init()
	Init()

	vc := `{"name": "Alice", "age": 25, "birth_date": "20000101", "edu_level": 4, "serial_no": "1234567890"}`
	var _vc zkclient.VC
	err := json.Unmarshal([]byte(vc), &_vc)
	assert.NoError(t, err)

	ZkUploadInput := ZkUploadInput{
		Vc:                 &_vc,
		BirthdateThreshold: "19990101",
	}
	key, iv := "verysecretkey123", "uniqueiv12345678"
	// zk := NewZk()

	t.Run("upload vc", func(_t *testing.T) {
		zkUploadOutput, err := NewZk().UploadVc(&ZkUploadInput, key, iv)
		assert.NoError(_t, err)
		fmt.Printf("zk upload output: %v\n", zkUploadOutput)
	})

	t.Run("zk proof", func(_t *testing.T) {
		flowProofForZk := FlowProofForZk{
			VcDataRoot: common.HexToHash("0xcca4353c87fe7c16438b81f6931e04112daa4c6c880e0df49a8950f27b4bcc23"),
			FlowRoot:   common.HexToHash("0x3efcbfbda4ea3ddd932e9c4959f54bbe5fc929359540a3405df28ccff746cec7"),
			Lemma: []common.Hash{
				common.HexToHash("0xcca4353c87fe7c16438b81f6931e04112daa4c6c880e0df49a8950f27b4bcc23"),
				common.HexToHash("0x9ac142c82258f2ce89d5d694b467d92caf1b4c97a33a34ddfcf1a78078b6cf18"),
				common.HexToHash("0xa1520264ae93cac619e22e8718fc4fa7ebdd23f493cad602434d2a58ff4868fb"),
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
			Path: 18945,
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
	})

	t.Run("zk verify", func(_t *testing.T) {
		proof := "d9ca38128b8e13b80ea6370705c7d66fedd483a9a8b565da4762ce079cd8572dc55f0f73389ad64708c91e326ff5f9fd7a3081ae3c1f4b535e9e7bda18398220808be6c1cb7b57154e8937d1c7785532907b99dacb1131818592af7fef65029f6d8130636c460b7efa461c8cbba67043127e768f7e5aaed69f6ead6de3afeb96"
		flowRoot := common.HexToHash("0x3efcbfbda4ea3ddd932e9c4959f54bbe5fc929359540a3405df28ccff746cec7")

		// proof := "d9ca38128b8e13b80ea6370705c7d66fedd483a9a8b565da4762ce079cd8572dc55f0f73389ad64708c91e326ff5f9fd7a3081ae3c1f4b535e9e7bda18398220808be6c1cb7b57154e8937d1c7785532907b99dacb1131818592af7fef65029f6d8130636c460b7efa461c8cbba67043127e768f7e5aaed69f6ead6de3afeb96"
		// flowRoot := common.HexToHash("0x3efcbfbda4ea3ddd932e9c4959f54bbe5fc929359540a3405df28ccff746cec7")

		isSucess, err := NewZk().ZkVerify(proof, "19990101", flowRoot.Hex())
		assert.NoError(t, err)
		assert.True(t, isSucess)
		fmt.Println(isSucess)
	})

}

// func TestUploadData(t *testing.T) {
// 	os.Chdir("/Users/dayong/myspace/mywork/storage-tool")
// 	config.Init()
// 	Init()

// 	// // data := make([]byte, 257)
// 	// vc := `{"name": "Alice", "age": 25, "birth_date": "20000101", "edu_level": 4, "serial_no": "1234567890"}`
// 	// var _vc zkclient.VC
// 	// err := json.Unmarshal([]byte(vc), &_vc)
// 	// assert.NoError(t, err)

// 	// data := _vc.Hash()
// 	// submissionTx, dataRoot, err := DefaultUploader().UploadBytes(data[:])
// 	// assert.NoError(t, err)

// 	// fmt.Printf("submission tx: %v\n", submissionTx)
// 	// fmt.Printf("data root: %v\n", dataRoot)

// 	vc := `{"name": "Alice", "age": 25, "birth_date": "20000101", "edu_level": 4, "serial_no": "1234567890"}`
// 	var _vc zkclient.VC
// 	err := json.Unmarshal([]byte(vc), &_vc)
// 	assert.NoError(t, err)

// 	ZkUploadInput := ZkUploadInput{
// 		Vc:                 &_vc,
// 		BirthdateThreshold: "19990101",
// 	}

// 	zk := NewZk()
// 	key, iv := "verysecretkey123", "uniqueiv12345678"
// 	zkUploadOutput, err := zk.UploadVc(&ZkUploadInput, key, iv)
// 	assert.NoError(t, err)

// 	fmt.Printf("zk upload output: %v\n", zkUploadOutput)
// }
// func TestZkProof(t *testing.T) {
// 	os.Chdir("/Users/dayong/myspace/mywork/storage-tool")
// 	config.Init()
// 	Init()

// 	proveInput := &ZkProofInput{
// 		ZkUploadInput:  ZkUploadInput,
// 		FlowProofForZk: zkUploadOutput.FlowProofForZk,
// 		Key:            key,
// 		IV:             iv,
// 	}
// 	proveOutput, err := zk.ZkProof(proveInput)
// 	assert.NoError(t, err)
// 	fmt.Printf("prove output: %v\n", proveOutput)
// }

func TestZkVerify(t *testing.T) {
	os.Chdir("/Users/dayong/myspace/mywork/storage-tool")
	config.Init()
	Init()

	proof := "179147f5ce659de5bb82c69649c2a296e9a3157a4d5a22696558af9d11e878875b3b11090e89240637b7c3e9a04bdb1663fbf2fb768ee08cde1f8b252214f5042217ab56966124501e54154e5822bb67720dbc82bcc2e8213bfa49336ec563086194563dfb03572b52cf75948f0e97ea1973a3d63330fd9b569db262c1c3bda1"
	flowRoot := common.HexToHash("0x096092f289fe8d67ff2ad798845fd059a82207025c0e6cd8b4fc344f30d53aa9")

	isSucess, err := NewZk().ZkVerify(proof, "20000101", flowRoot.Hex())
	assert.NoError(t, err)
	assert.True(t, isSucess)
	fmt.Println(isSucess)
}

// func TestGetChunksTree(t *testing.T) {
// 	data := make([]byte, 257)
// 	// vc := `{"name": "Alice", "age": 25, "birth_date": "20000101", "edu_level": 4, "serial_no": "1234567890"}`
// 	tree, err := DefaultUploader().getChunksTree(data)
// 	assert.NoError(t, err)

// 	fmt.Printf("root %v\n", tree.Root())
// 	fmt.Printf("proof at 0: %v\n", tree.ProofAt(0))
// }
