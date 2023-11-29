package utils

import (
	"fmt"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/grpc-boot/base/v2/logger"
)

func ParallelLoad(loadList ...func() error) {
	if len(loadList) < 1 {
		return
	}

	var wa sync.WaitGroup
	wa.Add(len(loadList))

	for _, load := range loadList {
		go func(l func() error) {
			st := time.Now()

			if err := l(); err != nil {
				logger.ZapError("init failed",
					logger.Error(err),
				)
			}

			pkg := runtime.FuncForPC(reflect.ValueOf(l).Pointer()).Name()
			fmt.Printf("load %v cost:%v\n", pkg, time.Since(st))

			wa.Done()
		}(load)
	}

	wa.Wait()
	return
}
