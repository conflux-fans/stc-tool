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

func TestZkProof(t *testing.T) {
	os.Chdir("/Users/dayong/myspace/mywork/storage-tool")
	config.Init()
	Init()

	vc := `{"name": "Alice", "age": 25, "birth_date": "20000101", "edu_level": 4, "serial_no": "1234567890"}`
	proveOutput, err := NewZk().ZkProof(vc, "verysecretkey123", "uniqueiv12345678", "19990101")
	assert.NoError(t, err)
	fmt.Printf("prove output: %v\n", proveOutput)
}

func TestGetChunksTree(t *testing.T) {
	data := make([]byte, 257)
	// vc := `{"name": "Alice", "age": 25, "birth_date": "20000101", "edu_level": 4, "serial_no": "1234567890"}`
	tree, err := DefaultUploader().getChunksTree(data)
	assert.NoError(t, err)

	fmt.Printf("root %v\n", tree.Root())
	fmt.Printf("proof at 0: %v\n", tree.ProofAt(0))
}

func TestUploadData(t *testing.T) {
	os.Chdir("/Users/dayong/myspace/mywork/storage-tool")
	config.Init()
	Init()

	// data := make([]byte, 257)
	vc := `{"name": "Alice", "age": 25, "birth_date": "20000101", "edu_level": 4, "serial_no": "1234567890"}`
	var _vc zkclient.VC
	err := json.Unmarshal([]byte(vc), &_vc)
	assert.NoError(t, err)

	data := _vc.Hash()
	submissionTx, dataRoot, err := DefaultUploader().UploadBytes(data[:])
	assert.NoError(t, err)

	fmt.Printf("submission tx: %v\n", submissionTx)
	fmt.Printf("data root: %v\n", dataRoot)
}

func TestZkVerify(t *testing.T) {
	os.Chdir("/Users/dayong/myspace/mywork/storage-tool")
	config.Init()
	Init()

	proof := "179147f5ce659de5bb82c69649c2a296e9a3157a4d5a22696558af9d11e878875b3b11090e89240637b7c3e9a04bdb1663fbf2fb768ee08cde1f8b252214f5042217ab56966124501e54154e5822bb67720dbc82bcc2e8213bfa49336ec563086194563dfb03572b52cf75948f0e97ea1973a3d63330fd9b569db262c1c3bda1"
	flowRoot := common.HexToHash("0x096092f289fe8d67ff2ad798845fd059a82207025c0e6cd8b4fc344f30d53aa9")

	isSucess, err := NewZk().ZkVerify(proof, "19990101", flowRoot.Hex())
	assert.NoError(t, err)
	fmt.Println(isSucess)
}
