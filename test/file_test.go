package test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenFile(t *testing.T) {
	// 设置测试文件名
	filename := "../tmp/test_file"
	fullFilename := filename + ".download"

	// 确保测试结束后清理文件
	defer os.Remove(fullFilename)

	// 测试文件创建和打开
	t.Run("创建并打开文件", func(t *testing.T) {
		file, err := os.OpenFile(fullFilename, os.O_RDWR|os.O_CREATE, 0666)
		assert.NoError(t, err, "打开文件时不应该出错")
		assert.NotNil(t, file, "文件对象不应为nil")

		// 如果成功打开，记得关闭文件
		if file != nil {
			file.Close()
		}

		// 验证文件是否真的被创建了
		_, err = os.Stat(fullFilename)
		assert.NoError(t, err, "文件应该已经被创建")
	})

	// 测试重复打开文件
	t.Run("重复打开已存在的文件", func(t *testing.T) {
		file, err := os.OpenFile(fullFilename, os.O_RDWR|os.O_CREATE, 0666)
		assert.NoError(t, err, "重新打开文件时不应该出错")
		assert.NotNil(t, file, "文件对象不应为nil")

		if file != nil {
			file.Close()
		}
	})
}
