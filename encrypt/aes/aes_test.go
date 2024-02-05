package aes

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"testing"

	"gotest.tools/assert"
)

func TestEncryptBytes(t *testing.T) {
	input := []byte("1")
	output := bytes.NewBuffer(make([]byte, 0))

	encryptor := new(AesEncryptor)
	err := encryptor.Encryptx(bytes.NewBuffer(input), output, []byte("1234567812345678"))
	assert.NilError(t, err)

	r, _ := io.ReadAll(output)
	fmt.Printf("result %x\n", r)
}

func TestDecryptBytes(t *testing.T) {
	input, err := hex.DecodeString("61626364656631323334353637383930cfab6d1815ef7a19aeba3b700c9d8c99")
	assert.NilError(t, err)
	fmt.Printf("input %x\n", input)

	inputBuf := bytes.NewBuffer(input)
	outputBuf := bytes.NewBuffer(make([]byte, 0))

	encryptor := new(AesEncryptor)
	err = encryptor.Decryptx(inputBuf, outputBuf, []byte("1234567812345678"))
	assert.NilError(t, err)

	r, _ := io.ReadAll(outputBuf)
	fmt.Printf("decrypt result %s\n", r)
}

func TestEncrypt(t *testing.T) {
	encryptor := new(AesEncryptor)
	encypted, err := encryptor.Encrypt([]byte("1"), []byte("1234567812345678"))
	assert.NilError(t, err)

	fmt.Printf("ecrypt result %x\n", encypted)
}

func TestDecrypt(t *testing.T) {
	input, err := hex.DecodeString("61626364656631323334353637383930cfab6d1815ef7a19aeba3b700c9d8c99")
	assert.NilError(t, err)
	fmt.Printf("input %x\n", input)

	encryptor := new(AesEncryptor)
	decypted, err := encryptor.Decrypt(input, []byte("1234567812345678"))
	assert.NilError(t, err)

	fmt.Printf("decrypt result %s\n", decypted)
}
