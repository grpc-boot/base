package utils

import (
	"os"
	"syscall"
)

func FileTime(fileName string) (createTime, lastAccessTime, lastWriteTime int64, err error) {
	info, err := os.Stat(fileName)
	if err != nil {
		return
	}

	attr := info.Sys().(*syscall.Win32FileAttributeData)
	createTime = attr.CreationTime.Nanoseconds() / 1e9
	lastAccessTime = attr.LastAccessTime.Nanoseconds() / 1e9
	lastWriteTime = attr.LastWriteTime.Nanoseconds() / 1e9
	return
}
