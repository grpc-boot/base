package utils

import (
	"os"

	"github.com/grpc-boot/base/v2/internal"
)

// FileExists 判断文件是否存在
func FileExists(fileName string) (exists bool, err error) {
	_, err = os.Stat(fileName)
	if err == nil {
		exists = true
		return
	}

	if os.IsNotExist(err) {
		err = nil
	}

	return
}

// FileTime 获取文件时间
func FileTime(fileName string) (createTime, lastAccessTime, lastWriteTime int64, err error) {
	return internal.FileTime(fileName)
}

// MkDir 创建目录
func MkDir(dir string, perm os.FileMode) (err error) {
	exists, err := FileExists(dir)
	if err != nil {
		return
	}

	if exists {
		return
	}

	return os.MkdirAll(dir, perm)
}
