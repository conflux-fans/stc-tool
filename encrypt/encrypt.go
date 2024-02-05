package encrypt

import (
	"errors"

	"github.com/zero-gravity-labs/zerog-storage-tool/encrypt/aes"
	"github.com/zero-gravity-labs/zerog-storage-tool/encrypt/empty"
)

type Encryptor interface {
	Encrypt(data, key []byte) ([]byte, error)
	Decrypt(ciphertext, key []byte) ([]byte, error)
}

var (
	aseEncryptor   aes.AesEncryptor
	emptyEncryptor empty.EmptyEncryptor
)

func GetEncryptor(method string) (Encryptor, error) {
	switch method {
	case "":
		return &emptyEncryptor, nil
	case "aes":
		return &aseEncryptor, nil
	}
	return nil, errors.New("unsupport")
}
