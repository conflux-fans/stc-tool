package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
	"io"
)

type AesEncryptor struct {
}

var iv = []byte("abcdef1234567890")

// // 使用AES加密
// func (a *AesEncryptor) Encrypt(data, key []byte) ([]byte, error) {
// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// 使用随机的IV向量
// 	ciphertext := make([]byte, aes.BlockSize+len(data))
// 	iv := ciphertext[:aes.BlockSize]
// 	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
// 		return nil, err
// 	}

// 	// 使用CBC模式加密
// 	mode := cipher.NewCBCEncrypter(block, iv)
// 	mode.CryptBlocks(ciphertext[aes.BlockSize:], data)

// 	return ciphertext, nil
// }

// // 使用AES解密
// func (a *AesEncryptor) Decrypt(ciphertext, key []byte) ([]byte, error) {
// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// 提取IV
// 	iv := ciphertext[:aes.BlockSize]
// 	ciphertext = ciphertext[aes.BlockSize:]

// 	// 使用CBC模式解密
// 	mode := cipher.NewCBCDecrypter(block, iv)
// 	mode.CryptBlocks(ciphertext, ciphertext)

// 	return ciphertext, nil
// }

func (a *AesEncryptor) Encryptx(input io.Reader, output io.Writer, key []byte) error {
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

		fmt.Println("read", n)

		paded := pad(buf[:n])
		cipher.NewCBCEncrypter(block, iv).CryptBlocks(paded, paded)
		n, err = output.Write(paded)
		if err != nil {
			return err
		}

		fmt.Println("write", n)
	}

	return nil
}

func (a *AesEncryptor) Decryptx(input io.Reader, output io.Writer, key []byte) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	iv := make([]byte, aes.BlockSize)
	_, err = input.Read(iv)
	if err != nil {
		return err
	}

	for {

		buf := make([]byte, 4096)
		n, err := input.Read(buf)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		fmt.Println("read", n)

		chunk := buf[:n]
		cipher.NewCBCDecrypter(block, iv).CryptBlocks(chunk, chunk)
		n, err = output.Write(chunk)
		if err != nil {
			return err
		}

		fmt.Println("write", n)
	}
	return nil
	// if len(ciphertext) < aes.BlockSize {
	// 	return nil, errors.New("ciphertext too short")
	// }
	// // iv := ciphertext[:aes.BlockSize]
	// fmt.Printf("decrypt iv %x\n", iv)
	// ciphertext = ciphertext[aes.BlockSize:]
	// cipher.NewCBCDecrypter(block, iv).CryptBlocks(ciphertext, ciphertext)
	// return ciphertext, nil
}

func (a *AesEncryptor) Encrypt(input, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	leng := len(input)
	if leng%16 != 0 {
		leng = leng/16*16 + 16
		leng = leng - len(input)
		for i := 0; i < leng; i++ {
			input = append(input, 0)
		}
		leng = len(input)
	}

	cipherText := make([]byte, aes.BlockSize+leng)
	copy(cipherText[:aes.BlockSize], iv)

	// if _, err := io.ReadFull(rand.Reader, iv); err != nil {
	// 	return nil, err
	// }
	fmt.Printf("encrypt iv %x\n", iv)
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(cipherText[aes.BlockSize:], input)
	return cipherText, nil
}

func (a *AesEncryptor) Decrypt(ciphertext, key []byte) ([]byte, error) {
	// func AesDecrpt(d string, key []byte) (string, error) {
	// ciphertext, err := hex.DecodeString(d)
	// if err != nil {
	// 	return "", err
	// }
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	// iv := ciphertext[:aes.BlockSize]
	fmt.Printf("decrypt iv %x\n", iv)
	ciphertext = ciphertext[aes.BlockSize:]
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(ciphertext, ciphertext)
	return ciphertext, nil
}

func pad(input []byte) []byte {
	leng := len(input)
	if leng%16 != 0 {
		leng = leng/16*16 + 16
		leng = leng - len(input)
		for i := 0; i < leng; i++ {
			input = append(input, 0)
		}
	}
	return input
}
