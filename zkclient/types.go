package zkclient

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/conflux-fans/storage-cli/encrypt/aes"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type VC struct {
	// {"name": "Alice", "age": 25, "birth_date": "20000101", "edu_level": 4, "serial_no": "1234567890"}
	Name      string `json:"name"`
	Age       byte   `json:"age"`
	BirthDate string `json:"birth_date"`
	EduLevel  byte   `json:"edu_level"`
	SerialNo  string `json:"serial_no"` // hex string
}

func (v *VC) Encode() []byte {
	bName := make([]byte, 16)
	fillBytes([]byte(v.Name), bName)

	birthDate := u64ToBytes(MustParseBirthDate(v.BirthDate))
	bBirth := make([]byte, 8)
	fillBytes(birthDate[:], bBirth)

	serial, _ := hex.DecodeString(v.SerialNo)
	bSerial := make([]byte, 32)
	fillBytes(serial, bSerial)

	encoded := make([]byte, 79)
	i := copy(encoded, []byte("name"))
	i += copy(encoded[i:], bName)
	i += copy(encoded[i:], []byte("age"))
	i += copy(encoded[i:], []byte{v.Age})
	i += copy(encoded[i:], []byte("birth"))
	i += copy(encoded[i:], bBirth)
	i += copy(encoded[i:], []byte("edu"))
	i += copy(encoded[i:], []byte{v.EduLevel})
	i += copy(encoded[i:], []byte("serial"))
	i += copy(encoded[i:], bSerial)

	if i != 79 {
		panic(fmt.Sprintf("valid bytes length should be 79, got %v", i))
	}

	return encoded
}

func (v *VC) EncodeAndPadToSector() []byte {
	return PadToSector(v.Encode())
}

func (v *VC) Hash() common.Hash {
	return crypto.Keccak256Hash(v.Encode())
}

func (v *VC) PlainText() []byte {
	return append(v.Encode(), v.Hash().Bytes()...)
}

func (v *VC) CipherText(key string, iv string) ([]byte, error) {
	encryptor := aes.NewAesCtrEncryptor(iv)
	buf := bytes.NewBuffer(nil)
	err := encryptor.Encrypt(bytes.NewReader(v.PlainText()), buf, []byte(key))
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (v *VC) MustGetUploadText(key, iv string) []byte {
	ct, err := v.CipherText(key, iv)
	if err != nil {
		panic(err)
	}
	data := append([]byte(iv), ct...)
	return PadToSector(data)
}

type ProveInput struct {
	Key         [16]byte          `json:"key"`
	Iv          [16]byte          `json:"iv"`
	Data        VC                `json:"data"`
	MerkleProof []common.Hash     `json:"merkle_proof"`
	PathIndex   uint64            `json:"path_index"`
	Extensions  []ExtensionSignal `json:"extensions"`
}

func NewProveInput(key string, iv string, data VC, merkleProof []common.Hash, pathIndex uint64, extensions []ExtensionSignal) *ProveInput {
	return &ProveInput{
		Key:         stringToByte16([]byte(key)),
		Iv:          stringToByte16([]byte(iv)),
		Data:        data,
		MerkleProof: merkleProof,
		PathIndex:   pathIndex,
		Extensions:  extensions,
	}
}

func (p *ProveInput) MarshalJSON() ([]byte, error) {
	type Alias ProveInput
	return json.Marshal(&struct {
		Key string `json:"key"`
		Iv  string `json:"iv"`
		*Alias
	}{
		Key:   hex.EncodeToString(p.Key[:]),
		Iv:    hex.EncodeToString(p.Iv[:]),
		Alias: (*Alias)(p),
	})
}

func (p *ProveInput) UnmarshalJSON(data []byte) error {
	type Alias ProveInput
	aux := &struct {
		Key string `json:"key"`
		Iv  string `json:"iv"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	key, err := hex.DecodeString(aux.Key)
	if err != nil {
		return err
	}
	iv, err := hex.DecodeString(aux.Iv)
	if err != nil {
		return err
	}
	copy(p.Key[:], key)
	copy(p.Iv[:], iv)
	return nil
}

type ExtensionSignal struct {
	Date   *time.Time `json:"date,omitempty"`
	Number *uint64    `json:"number,omitempty"`
}

func (e ExtensionSignal) MarshalJSON() ([]byte, error) {
	type Alias ExtensionSignal
	aux := &struct {
		Date string `json:"date,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(&e),
	}
	if e.Date != nil {
		aux.Date = e.Date.Format("20060102")
	}
	return json.Marshal(aux)
}

func (e *ExtensionSignal) UnmarshalJSON(data []byte) error {
	type Alias ExtensionSignal
	aux := &struct {
		Date string `json:"date,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(e),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.Date != "" {
		date, err := time.Parse("20060102", aux.Date)
		if err != nil {
			return err
		}
		e.Date = &date
	}
	return nil
}

type ProveOutput struct {
	Proof            string
	VcUploadTextRoot common.Hash
	FlowRoot         common.Hash
}

type VerifyInput struct {
	Extensions []ExtensionSignal `json:"extensions"`
	Root       common.Hash       `json:"root"`
}
