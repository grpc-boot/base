package base

import (
	"golang.org/x/sys/windows"
	"syscall"
)

var modkernel32 = windows.NewLazySystemDLL("kernel32.dll")
var procSetProcessAffinityMask = modkernel32.NewProc("SetProcessAffinityMask")

// workWithAffinity 绑定CPU
func workWithAffinity(id int) (uint64, error) {
	taskHandle := windows.CurrentProcess()
	err := windows.SetPriorityClass(taskHandle, windows.HIGH_PRIORITY_CLASS)
	if err != nil {
		return 0, err
	}

	err = setProcessAffinityMask(taskHandle, 1<<id)
	if err != nil {
		return 0, err
	}

	return uint64(windows.GetCurrentThreadId()), nil
}

func setProcessAffinityMask(handle windows.Handle, processAffinityMask uint64) (err error) {
	r1, _, e1 := syscall.Syscall(procSetProcessAffinityMask.Addr(), 2, uintptr(handle), uintptr(processAffinityMask), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = e1
		} else {
			err = syscall.EINVAL
		}
	}
	return
}
