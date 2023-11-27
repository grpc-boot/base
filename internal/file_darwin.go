package internal

import (
	"os"
	"syscall"
)

func FileTime(fileName string) (createTime, lastAccessTime, lastWriteTime int64, err error) {
	info, err := os.Stat(fileName)
	if err != nil {
		return
	}

	attr := info.Sys().(*syscall.Stat_t)
	createTime = attr.Ctimespec.Sec
	lastAccessTime = attr.Atimespec.Sec
	lastWriteTime = attr.Mtimespec.Sec
	return
}
