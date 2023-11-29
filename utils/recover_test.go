package utils

import "testing"

func TestRecover(t *testing.T) {
	go Recover("recover test", func() {
		panic("panic with test")
	})
}
