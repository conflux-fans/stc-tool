package encrypt

import (
	"bytes"

	"io"
	"os"

	"github.com/conflux-fans/storage-cli/constants/enums"
	"github.com/conflux-fans/storage-cli/encrypt/aes"
	"github.com/conflux-fans/storage-cli/encrypt/empty"
	"github.com/pkg/errors"
)

type Encryptor interface {
	Encrypt(input io.Reader, output io.Writer, key []byte) error
	Decrypt(input io.Reader, output io.Writer, key []byte) error
}

var (
	aseCbcEncryptor aes.AesCbcEncryptor
	aseCtrEncryptor aes.AesCtrEncryptor = *aes.NewAesCtrEncryptor(string(aes.IV))
	emptyEncryptor  empty.EmptyEncryptor
)

func GetEncryptor(method enums.CipherMethod) (Encryptor, error) {
	switch method {
	case enums.CIPHER_METHOD_EMPTY:
		return &emptyEncryptor, nil
	case enums.CIPHER_METHOD_AES_CBC:
		return &aseCbcEncryptor, nil
	case enums.CIPHER_METHOD_AES_CTR:
		return &aseCtrEncryptor, nil
	}
	return nil, errors.New("unsupport cipher method")
}

func EncryptBytes(e Encryptor, input, key []byte) ([]byte, error) {
	inputBuf := bytes.NewBuffer(input)
	outputBuf := bytes.NewBuffer(make([]byte, 0))

	if err := e.Encrypt(inputBuf, outputBuf, key); err != nil {
		return nil, err
	}
	return io.ReadAll(outputBuf)
}

func DecryptBytes(e Encryptor, input, key []byte) ([]byte, error) {
	inputBuf := bytes.NewBuffer(input)
	outputBuf := bytes.NewBuffer(make([]byte, 0))

	if err := e.Decrypt(inputBuf, outputBuf, key); err != nil {
		return nil, err
	}
	return io.ReadAll(outputBuf)
}

func EncryptFile(e Encryptor, source, outputDirPath string, key []byte) (string, error) {
	if err := os.MkdirAll(outputDirPath, 0755); err != nil {
		return "", errors.WithMessage(err, "Failed to create directory")
	}

	sf, err := os.OpenFile(source, os.O_RDONLY, 0666)
	if err != nil {
		return "", errors.WithMessage(err, "Failed to open source file")
	}
	defer sf.Close()

	outputhPath := outputDirPath + mustGetFileName(sf) + ".encrypt"

	of, err := os.OpenFile(outputhPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return "", errors.WithMessage(err, "Failed to create output file")
	}
	defer of.Close()

	return outputhPath, e.Encrypt(sf, of, key)
}

func DecryptFile(e Encryptor, source, outputDirPath string, key []byte) (string, error) {
	// fmt.Printf("decrypt file source %s, out %s\n", source, outputDirPath)
	if err := os.MkdirAll(outputDirPath, 0755); err != nil {
		return "", errors.WithMessage(err, "Failed to create directory")
	}

	sf, err := os.OpenFile(source, os.O_RDONLY, 0666)
	if err != nil {
		return "", errors.WithMessage(err, "Failed to open source file")
	}
	defer sf.Close()
	// fmt.Printf("sf name %s\n", mustGetFileName(sf))

	outputhPath := outputDirPath + mustGetFileName(sf) + ".decrypt"

	of, err := os.OpenFile(outputhPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return "", errors.WithMessage(err, "Failed to create output file")
	}
	defer of.Close()

	return outputhPath, e.Decrypt(sf, of, key)
}

func mustGetFileName(f *os.File) string {
	stat, err := f.Stat()
	if err != nil {
		panic(err)
	}
	return stat.Name()
}
