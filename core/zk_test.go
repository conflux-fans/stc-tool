package core

import (
	"encoding/hex"
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
	_, err := ZkProof(vc, "20240101")
	assert.NoError(t, err)
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
	segTree, trunksTree, err := DefaultUploader().UploadString(data[:])
	assert.NoError(t, err)

	fmt.Printf("seg tree root: %v\n", segTree.Root())
	fmt.Printf("seg proof at 0: %v\n", segTree.ProofAt(0))
	fmt.Printf("trunks tree root: %v\n", trunksTree.Root())
	fmt.Printf("trunks proof at 0: %v\n", trunksTree.ProofAt(0))

	// FIXME: get flow root
	// storageSystemRoot, err := defaultFlow.Root(nil)
	storageSystemRoot := common.Hash{}
	assert.NoError(t, err)
	fmt.Printf("storage system root: %x\n", storageSystemRoot)
}

func TestZkVerify(t *testing.T) {
	os.Chdir("/Users/dayong/myspace/mywork/storage-tool")
	config.Init()
	Init()

	proof := "dcd19771587f25cb7d706020efb93cae9b8898116932074fc389311b035a5b967d564463c426aaa34960f5b50859f298a7e62d9fac74c5cc680ae7365597d01a79673058eb3227b58706d7d6c92d41e6146fdb45442ee2ddc622315d380cdca449b54240ae7fa68e037f58092031f23da1e99aae84df13a400ea961bf73f798f"

	// FIXME: get flow root
	// storageSystemRoot, err := defaultFlow.Root(nil)
	// assert.NoError(t, err)
	storageSystemRoot := common.Hash{}

	isSucess, err := ZkVerify(proof, "20240101", hex.EncodeToString(storageSystemRoot[:]))
	assert.NoError(t, err)
	fmt.Println(isSucess)
}
