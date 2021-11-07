package base

import (
	"fmt"
	"os"
)

func Black(format string, args ...interface{}) {
	fmt.Printf("\033[0;30m"+format+"\033[0m\n", args...)
}

func Red(format string, args ...interface{}) {
	fmt.Printf("\033[0;31m"+format+"\033[0m\n", args...)
}

func RedFatal(format string, args ...interface{}) {
	Red(format, args...)
	os.Exit(1)
}

func Green(format string, args ...interface{}) {
	fmt.Printf("\033[0;32m"+format+"\033[0m\n", args...)
}

func Blue(format string, args ...interface{}) {
	fmt.Printf("\033[0;34m"+format+"\033[0m\n", args...)
}

func Yellow(format string, args ...interface{}) {
	fmt.Printf("\033[0;33m"+format+"\033[0m\n", args...)
}

func Fuchsia(format string, args ...interface{}) {
	fmt.Printf("\033[0;35m"+format+"\033[0m\n", args...)
}
