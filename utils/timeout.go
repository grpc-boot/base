package utils

import (
	"context"
	"time"
)

func Timeout(timeout time.Duration, handler func()) error {
	var (
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
		done        = make(chan struct{}, 1)
	)
	defer cancel()

	go func() {
		handler()
		done <- struct{}{}
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
