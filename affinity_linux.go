package base

import (
	"golang.org/x/sys/unix"
	"syscall"
)

// workWithAffinity 绑定CPU
func workWithAffinity(id int) (uint64, error) {
	var newMask unix.CPUSet
	newMask.Set(id)

	tid := syscall.Gettid()

	err := unix.SchedSetaffinity(tid, &newMask)
	if err != nil {
		return 0, err
	}

	return uint64(tid), nil
}
