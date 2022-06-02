package base

import (
	"os"
	"path/filepath"
)

var (
	binPath string
)

func init() {
	binPath = filepath.Dir(os.Args[0])
}

// BinPath 获取Bin目录
func BinPath() string {
	return binPath
}
