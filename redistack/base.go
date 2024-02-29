package redistack

import (
	"context"
	"time"

	"github.com/redis/rueidis"
)

func Do(ctx context.Context, client rueidis.Client, cmd rueidis.Completed) rueidis.RedisResult {
	return client.Do(ctx, cmd)
}

func DoWithTimeout(timeout time.Duration, client rueidis.Client, cmd rueidis.Completed) rueidis.RedisResult {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return client.Do(ctx, cmd)
}
