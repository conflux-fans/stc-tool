package logger

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

var _logger *logrus.Logger

func init() {
	// 创建一个新的 Logger 实例
	_logger = logrus.New()
	// 创建一个自定义的 Hook，并添加到 Logger 中
	_logger.AddHook(&prefixHook{Prefix: "\x1b[42m[TOOL]\x1b[0m"})
}

func Fail(desc string) {
	fmt.Println("\n❌ \x1b[31mFAIL\x1b[0m: " + desc + "\n")
}

func Failf(descFormat string, a ...any) {
	fmt.Printf("\n❌ \x1b[31mFAIL\x1b[0m: "+descFormat+"\n", a...)
}

func FailfWithParams(params map[string]string, descFormat string, a ...any) {
	result := resultByParams(params)
	fmt.Printf("\n❌ \x1b[31mFAIL\x1b[0m: == "+descFormat+" ==\n", a...)
	fmt.Println(result)
}

func Success(desc string) {
	fmt.Println("\n✅ \x1b[32mSUCCESS\x1b[0m: == " + desc + " ==\n")
}

func Successf(descFormat string, a ...any) {
	fmt.Printf("\n✅ \x1b[32mSUCCESS\x1b[0m: == "+descFormat+" ==\n", a...)
}

func SuccessWithResult(result string, descFormat string, a ...any) {
	fmt.Printf("\n✅ \x1b[32mSUCCESS\x1b[0m: == "+descFormat+" ==\n", a...)
	fmt.Println(result)
}

func SuccessfWithParams(params map[string]string, descFormat string, a ...any) {
	result := resultByParams(params)
	SuccessWithResult(result, descFormat, a...)
}

func SuccessfWithList[T any](list []T, descFormat string, a ...any) {
	// fmt.Printf("\n✅ \x1b[32mSUCCESS\x1b[0m: == "+descFormat+" ==\n", a...)
	var result string
	for _, v := range list {
		result += fmt.Sprintf("    - %v\n", v)
	}
	SuccessWithResult(result, descFormat, a...)
}

func resultByParams(params map[string]string) string {
	// 根据 key 最长的值的长度为对齐长度，创建一个分行的字符串
	maxKeyLen := 0
	for k := range params {
		if len(k) > maxKeyLen {
			maxKeyLen = len(k)
		}
	}

	var result string
	for k, v := range params {
		// result += fmt.Sprint(" -- %-"+fmt.Sprintf("%d", maxKeyLen)+"s: %s\n", k, v)
		result += fmt.Sprintf("    - %-*s: %s\n", maxKeyLen, k, v)
	}
	return result
}
func Get() *logrus.Logger {
	return _logger
}

// prefixHook 自定义的 Hook，用于在日志条目前添加前缀
type prefixHook struct {
	Prefix string
}

// Fire 实现 Hook 接口的 Fire 方法
func (hook *prefixHook) Fire(entry *logrus.Entry) error {
	// 在日志条目消息前添加前缀
	entry.Message = hook.Prefix + " " + entry.Message
	return nil
}

// Levels 实现 Hook 接口的 Levels 方法
func (hook *prefixHook) Levels() []logrus.Level {
	// 设置 Hook 对所有日志级别生效
	return logrus.AllLevels
}
