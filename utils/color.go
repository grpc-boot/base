package utils

import (
	"fmt"
	"os"
)

// Black 默认输出
func Black(format string, args ...any) {
	fmt.Printf("\033[0;30m"+format+"\033[0m\n", args...)
}

// Red 错误输出
func Red(format string, args ...any) {
	fmt.Printf("\033[0;31m"+format+"\033[0m\n", args...)
}

// RedFatal 致命错误输出
func RedFatal(format string, args ...any) {
	Red(format, args...)
	os.Exit(1)
}

// Green 成功输出
func Green(format string, args ...any) {
	fmt.Printf("\033[0;32m"+format+"\033[0m\n", args...)
}

// Blue 颜色输出
func Blue(format string, args ...any) {
	fmt.Printf("\033[0;34m"+format+"\033[0m\n", args...)
}

// Yellow 警告输出
func Yellow(format string, args ...any) {
	fmt.Printf("\033[0;33m"+format+"\033[0m\n", args...)
}

// Fuchsia 颜色输出
func Fuchsia(format string, args ...any) {
	fmt.Printf("\033[0;35m"+format+"\033[0m\n", args...)
}
