package gored

import (
	"errors"

	"github.com/grpc-boot/base/v2/logger"
	"github.com/redis/go-redis/v9"
)

type Cmd interface {
	Err() error
}

func DealCmdErr(cmd Cmd) error {
	err := cmd.Err()

	if errors.Is(err, redis.Nil) {
		err = nil
	}

	if err != nil {
		logger.ZapError("exec redis cmd failed",
			logger.Error(err),
		)
	}

	return err
}
