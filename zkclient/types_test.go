package zkclient

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVc(t *testing.T) {
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
		plainText      string
		cipherText     string
		cipherTextRoot string
	}{
		encoded:    "6e616d65416c696365000000000000000000000061676519626972746880436d38000000006564750473657269616c1234567890000000000000000000000000000000000000000000000000000000",
		vcHash:     "0xc941d9c835cd182571dafe1f7001cf98533103d977010dd6cc48e7bfafe525c1",
		plainText:  "6e616d65416c696365000000000000000000000061676519626972746880436d38000000006564750473657269616c1234567890000000000000000000000000000000000000000000000000000000c941d9c835cd182571dafe1f7001cf98533103d977010dd6cc48e7bfafe525c1",
		cipherText: "8eb770d23b1a772b947eda3b7e6b61d527ae7dfe4f05f35f1392e29d70a9ecfdf6e134f552ee7b22dd8f604247c8f2c2279f45b7b16edabec4ecf59226232fee1e565617715342029384d87b75e00fa87f6c35322be0ddb29e907872634ec963442b1436d7fc8fed431c66c19f53d2",
	}

	encodedVc := vc.Encode()
	cipherText, err := vc.CipherText("verysecretkey123", "uniqueiv12345678")
	assert.NoError(t, err)

	assert.Equal(t, expect.encoded, hex.EncodeToString(encodedVc))
	assert.Equal(t, expect.vcHash, vc.Hash().String())
	assert.Equal(t, expect.plainText, hex.EncodeToString(vc.PlainText()))
	assert.Equal(t, expect.cipherText, hex.EncodeToString(cipherText))
}
