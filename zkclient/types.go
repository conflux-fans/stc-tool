package zkclient

import (
	"encoding/hex"
	"fmt"

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

type ProveInput struct {
	Data               VC            `json:"data"`
	BirthdateThreshold string        `json:"birthdate_threshold"`
	MerkleProof        []common.Hash `json:"merkle_proof"`
	PathIndex          uint64        `json:"path_index"`
}

type ProveOutput struct {
	Proof         string
	EncryptVcRoot common.Hash
	FlowRoot      common.Hash
}

type VerifyInput struct {
	BirthdateThreshold string      `json:"birthdate_threshold"`
	Root               common.Hash `json:"root"`
}
