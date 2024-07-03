package zkclient

import (
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

func TestVcEncodeAndHash(t *testing.T) {
	vc := VC{
		Name:      "Alice",
		Age:       25,
		BirthDate: "20000101",
		EduLevel:  4,
		SerialNo:  "1234567890",
	}

	expect := struct {
		encoded        string
		vcHash         string
		vcRootInDepth0 string
	}{
		encoded:        "6e616d65416c696365000000000000000000000061676519626972746880436d38000000006564750473657269616c1234567890000000000000000000000000000000000000000000000000000000",
		vcHash:         "0xc941d9c835cd182571dafe1f7001cf98533103d977010dd6cc48e7bfafe525c1",
		vcRootInDepth0: "0x1efee730fb75c1fe4aa02f6af6b058248224af23952fc3c0380ccc5d25af6ef8",
	}

	encodedVc := vc.Encode()
	paddedVcHash := PadToSector(vc.Hash().Bytes())
	vc_root_in_depth_0 := crypto.Keccak256Hash(paddedVcHash)

	assert.Equal(t, expect.encoded, hex.EncodeToString(encodedVc))
	assert.Equal(t, expect.vcHash, vc.Hash().String())
	assert.Equal(t, expect.vcRootInDepth0, vc_root_in_depth_0.String())
}
