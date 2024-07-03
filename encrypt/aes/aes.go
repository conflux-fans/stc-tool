package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"io"

	"github.com/conflux-fans/storage-cli/logger"
)

type AesEncryptor struct {
}

var iv = []byte("abcdef1234567890")

func (a *AesEncryptor) Encrypt(input io.Reader, output io.Writer, key []byte) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	output.Write(iv)

	for {
		buf := make([]byte, 4096)
		n, err := input.Read(buf)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		// fmt.Println("read", n)

		paded := pad(buf[:n])
		cipher.NewCBCEncrypter(block, iv).CryptBlocks(paded, paded)
		n, err = output.Write(paded)
		if err != nil {
			return err
		}

		// fmt.Println("write", n)
	}

	return nil
}

func (a *AesEncryptor) Decrypt(input io.Reader, output io.Writer, key []byte) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	iv := make([]byte, aes.BlockSize)
	_, err = input.Read(iv)
	if err != nil {
		return err
	}

	end := false
	for {
		if end {
			break
		}

		buf := make([]byte, 4096)
		n, err := input.Read(buf)
		if err == io.EOF || n < 4096 {
			end = true
		} else if err != nil {
			return err
		}

		logger.Get().WithField("size", n).WithField("end", end).Info("read data")

		chunk := buf[:n]
		cipher.NewCBCDecrypter(block, iv).CryptBlocks(chunk, chunk)

		if end {
			chunk = trimTailZeros(chunk)
		}

		n, err = output.Write(chunk)
		if err != nil {
			return err
		}

		// fmt.Println("write", n)
		logger.Get().WithField("size", n).Info("write to output")
	}
	return nil
}

func pad(input []byte) []byte {
	padLen := aes.BlockSize - len(input)%aes.BlockSize
	input = append(input, make([]byte, padLen)...)
	return input
}

func trimTailZeros(input []byte) []byte {
	lastNonZeroIndex := len(input) - 1
	for lastNonZeroIndex >= 0 && input[lastNonZeroIndex] == 0 {
		lastNonZeroIndex--
	}

	result := input[:lastNonZeroIndex+1]
	return result
}
