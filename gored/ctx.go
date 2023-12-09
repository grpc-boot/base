package gored

import (
	"context"
	"time"
)

func TimeoutDo(timeout time.Duration, handler func(ctx context.Context)) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	handler(ctx)
}
