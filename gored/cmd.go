package gored

import (
	"github.com/grpc-boot/base/v2/logger"
	"github.com/redis/go-redis/v9"
)

type Cmd interface {
	Err() error
}

func DealCmdErr(cmd Cmd) error {
	err := cmd.Err()

	if err == redis.Nil {
		err = nil
	}

	if err != nil {
		logger.ZapError("exec redis cmd failed",
			logger.Error(err),
		)
	}

	return err
}
