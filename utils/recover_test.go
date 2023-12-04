package utils

import "testing"

func TestRecover(t *testing.T) {
	go Recover("recover test", func(args ...any) {
		panic("panic with test")
	})
}
