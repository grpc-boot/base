package base

import (
	"runtime"
)

// WorkWithAffinity 绑定CPU
func WorkWithAffinity(id int) (uint64, error) {
	runtime.LockOSThread()
	return workWithAffinity(id)
}
