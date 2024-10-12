package encryptutils

import (
	"github.com/conflux-fans/storage-cli/constants/enums"
	"github.com/conflux-fans/storage-cli/encrypt"
	"github.com/pkg/errors"
)

func EncryptFile(filePath string, method enums.CipherMethod, password string) (outputPath string, err error) {
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

func EncryptBytes(soruce []byte, method enums.CipherMethod, password string) ([]byte, error) {
	encryptor, err := encrypt.GetEncryptor(method)
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to get encryptor")
	}

	output, err := encrypt.EncryptBytes(encryptor, soruce, []byte(password))
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to encrypt")
	}
	return output, nil
}
