package enums

import "github.com/conflux-fans/storage-cli/pkg/utils/enumutils"

type CipherMethod int

const (
	CIPHER_METHOD_EMPTY CipherMethod = iota + 1
	CIPHER_METHOD_AES_CBC
	CIPHER_METHOD_AES_CTR
)

var cipherMethodEb enumutils.EnumBase[CipherMethod]

func init() {
	cipherMethodEb = enumutils.NewEnumBase("CIPHER_METHOD", map[CipherMethod]string{
		CIPHER_METHOD_EMPTY:   "",
		CIPHER_METHOD_AES_CBC: "AES_CBC",
		CIPHER_METHOD_AES_CTR: "AES_CTR",
	})
}

func (m CipherMethod) MarshalText() ([]byte, error) {
	return cipherMethodEb.MarshalText(m)
}

func (m *CipherMethod) UnmarshalText(data []byte) error {
	val, err := cipherMethodEb.UnmarshalText(data)
	if err != nil {
		return err
	}
	*m = val
	return nil
}

func (m CipherMethod) String() string {
	return cipherMethodEb.String(m)
}

func ParseCipherMethod(s string) (CipherMethod, error) {
	return cipherMethodEb.Parse(s)
}
