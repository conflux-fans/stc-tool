package test

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func generateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func createFolderWithRandomFiles(folderName string, numFiles int) error {
	err := os.MkdirAll(folderName, 0755)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("创建文件夹失败: %v", err)
	}

	for i := 0; i < numFiles; i++ {
		fileName := filepath.Join(folderName, fmt.Sprintf("file_%d.txt", i+1))
		contentLength := rand.Intn(100) + 1 // 1到100之间的随机数
		content := generateRandomString(contentLength)

		err := os.WriteFile(fileName, []byte(content), 0644)
		if err != nil {
			return fmt.Errorf("写入文件 %s 失败: %v", fileName, err)
		}
	}

	return nil
}

func TestGenFiles(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Errorf("获取用户主目录失败: %v", err)
		return
	}

	folderName := filepath.Join(homeDir, "tmp", "random_files")
	numFiles := 10 // 您可以更改这个数字来生成不同数量的文件

	err = createFolderWithRandomFiles(folderName, numFiles)
	if err != nil {
		t.Errorf("创建随机文件失败: %v", err)
		return
	}

	t.Logf("已在 '%s' 文件夹中创建 %d 个随机内容文件。", folderName, numFiles)
}
