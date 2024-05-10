package logger

import (
	"fmt"
	"testing"
)

func TestLogFail(t *testing.T) {
	Fail("Use Fail to print something")
	Failf("Use Failf to print something")
	Failf("Use Failf to print something %s", "arg1")
	FailfWithParams(map[string]string{"Foo": "bar", "Fooxxx": "zoo"}, "User FailfWithParams to print something %s", "arg1")
}

func TestLogSuccess(t *testing.T) {
	Success("Use Success to print something")
	Successf("Use Successf to print something")
	Successf("Use Successf to print something %s", "arg1")
	SuccessWithResult("Excute result", "Use SuccessWithResult to print something %s", "arg1")
	SuccessfWithParams(map[string]string{"Foo": "bar", "Fooxxx": "zoo"}, "User SuccessfWithParams to print something %s", "arg1")
}

func TestResultByParams(t *testing.T) {
	r := resultByParams(map[string]string{"Fooxxx": "bar", "For": "zoo"})
	fmt.Println(r)
}

func TestAlign(t *testing.T) {

	align := "-" // 左对齐使用 "-"，右对齐使用空字符串 ""
	minWidth := 10
	text := "dynamic"

	// 构建格式化字符串
	var format string
	if align == "-" {
		format = fmt.Sprintf("|%-*s|\n", minWidth, text) // 左对齐
	} else {
		format = fmt.Sprintf("|%*s|\n", minWidth, text) // 右对齐
	}

	// 输出格式化后的字符串
	fmt.Printf(format)

}
