package encryptutils

import (
	"github.com/pkg/errors"
	"github.com/zero-gravity-labs/zerog-storage-tool/encrypt"
)

func EncryptFile(filePath string, method string, password string) (outputPath string, err error) {
	encryptor, err := encrypt.GetEncryptor(method)
	if err != nil {
		return "", errors.WithMessage(err, "Failed to get encryptor")
	}

	outputPath, err = encrypt.EncryptFile(encryptor, filePath, "./", []byte(password))
	if err != nil {
		return "", errors.WithMessage(err, "Failed to encrypt")
	}
	return outputPath, nil
}
