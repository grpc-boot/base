package utils

import (
	"os"

	"github.com/grpc-boot/base/v2/internal"
)

// FileExists 判断文件是否存在
func FileExists(fileName string) bool {
	_, err := os.Stat(fileName)
	if err == nil {
		return true
	}

	return !os.IsNotExist(err)
}

// FileTime 获取文件时间
func FileTime(fileName string) (createTime, lastAccessTime, lastWriteTime int64, err error) {
	return internal.FileTime(fileName)
}
