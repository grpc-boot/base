package utils

import (
	"strings"
	"testing"
)

func TestJoin(t *testing.T) {
	ss := []string{"s1", "s2"}

	res1 := strings.Join(ss, ",")
	t.Logf("res1: %s", res1)

	is := []int{1, 2, 45}
	resInt := Join(",", is...)
	t.Logf("resInt: %s", resInt)

	i32s := []int32{1, 2, 45}
	resInt32 := Join(",", i32s...)
	t.Logf("resInt32: %s", resInt32)
}
