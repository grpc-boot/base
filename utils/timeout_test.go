package utils

import (
	"context"
	"testing"
	"time"
)

func TestTimeout(t *testing.T) {
	err := Timeout(time.Second, func() {
		time.Sleep(time.Millisecond * 500)
	})

	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}

	err = Timeout(time.Millisecond*100, func() {
		time.Sleep(time.Millisecond * 200)
	})

	if err != context.DeadlineExceeded {
		t.Fatalf("want err, got %v", err)
	}
}
