package utils

import (
	"os"

	"github.com/grpc-boot/base/v3/internal"
)

// FileExists 判断文件是否存在
func FileExists(fileName string) (exists bool, err error) {
	return internal.FileExists(fileName)
}

// FileTime 获取文件时间
func FileTime(fileName string) (createTime, lastAccessTime, lastWriteTime int64, err error) {
	return internal.FileTime(fileName)
}

// MkDir 创建目录
func MkDir(dir string, perm os.FileMode) (err error) {
	return internal.MkDir(dir, perm)
}
