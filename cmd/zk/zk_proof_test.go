package zk

import (
	"os"
	"testing"
)

func TestReadZkProofInput(t *testing.T) {
	// 创建临时的 zk_proof_input.json 文件
	tempContent := `{
		"Vc": {"name": "Alice", "age": 25, "birth_date": "20000101", "edu_level": 4, "serial_no": "1234567890"},
		"Key": "verysecretkey123",
		"IV": "uniqueiv12345678",
		"Birthdate_threshold": "2024-01-01"
	}`
	err := os.WriteFile("./zk_proof_input.json", []byte(tempContent), 0644)
	if err != nil {
		t.Fatalf("无法创建临时文件: %v", err)
	}
	defer os.Remove("./zk_proof_input.json") // 测试结束后删除临时文件

	input, err := readZkProofInput("./zk_proof_input.json")
	if err != nil {
		t.Fatalf("Failed to get zk proof input: %v", err)
	}
	t.Logf("input: %+v", input)
}
