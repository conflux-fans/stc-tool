package encrypt

import (
	"bytes"

	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/zero-gravity-labs/zerog-storage-tool/encrypt/aes"
	"github.com/zero-gravity-labs/zerog-storage-tool/encrypt/empty"
)

type Encryptor interface {
	Encrypt(input io.Reader, output io.Writer, key []byte) error
	Decrypt(input io.Reader, output io.Writer, key []byte) error
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

func EncryptFile(e Encryptor, source, outputDirPath string, key []byte) error {
	if err := os.MkdirAll(outputDirPath, 0755); err != nil {
		return errors.WithMessage(err, "Failed to create directory")
	}

	sf, err := os.OpenFile(source, os.O_RDONLY, 0666)
	if err != nil {
		return errors.WithMessage(err, "Failed to open source file")
	}

	of, err := os.OpenFile(outputDirPath+mustGetFileName(sf), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return errors.WithMessage(err, "Failed to create output file")
	}

	return e.Encrypt(sf, of, key)
}

func DecryptFile(e Encryptor, source, outputDirPath string, key []byte) error {
	// fmt.Printf("decrypt file source %s, out %s\n", source, outputDirPath)
	if err := os.MkdirAll(outputDirPath, 0755); err != nil {
		return errors.WithMessage(err, "Failed to create directory")
	}

	sf, err := os.OpenFile(source, os.O_RDONLY, 0666)
	if err != nil {
		return errors.WithMessage(err, "Failed to open source file")
	}
	// fmt.Printf("sf name %s\n", mustGetFileName(sf))

	of, err := os.OpenFile(outputDirPath+mustGetFileName(sf), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return errors.WithMessage(err, "Failed to create output file")
	}

	return e.Decrypt(sf, of, key)
}

func mustGetFileName(f *os.File) string {
	stat, err := f.Stat()
	if err != nil {
		panic(err)
	}
	return stat.Name()
}
