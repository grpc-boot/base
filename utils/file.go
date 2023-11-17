package utils

import (
	"os"
)

// FileExists 判断文件是否存在
func FileExists(fileName string) bool {
	_, err := os.Stat(fileName)
	if err == nil {
		return true
	}

	return !os.IsNotExist(err)
}
