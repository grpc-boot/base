package utils

import (
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/grpc-boot/base/v2/logger"

	"go.uber.org/zap"
)

func ParallelLoad(loadList ...func() error) {
	if len(loadList) < 1 {
		return
	}

	wa := &sync.WaitGroup{}
	wa.Add(len(loadList))

	for _, load := range loadList {
		go func(l func() error) {
			defer wa.Done()

			st := time.Now()
			if err := l(); err != nil {
				logger.Error(
					"init failed",
					zap.NamedError("Error", err),
				)
			} else {
				pkg := runtime.FuncForPC(reflect.ValueOf(l).Pointer()).Name()
				logger.Info("load component success",
					zap.String("Pkg", pkg),
					zap.Duration("Duration", time.Since(st)),
				)
			}
		}(load)
	}

	wa.Wait()
	return
}
