package empty

type EmptyEncryptor struct {
}

func (a *EmptyEncryptor) Encrypt(input, key []byte) ([]byte, error) {
	return input, nil
}

func (a *EmptyEncryptor) Decrypt(ciphertext, key []byte) ([]byte, error) {
	return ciphertext, nil
}
