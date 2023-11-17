package components

import "github.com/grpc-boot/base/core/gopool"

// NewGoPool 实例化goroutine池
func NewGoPool(size int, opts ...gopool.Option) (*gopool.Pool, error) {
	return gopool.NewPool(size, opts...)
}
