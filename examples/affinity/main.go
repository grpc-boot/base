package main

import (
	"time"

	"github.com/grpc-boot/base"
)

func main() {
	tId, err := base.WorkWithAffinity(0)

	if err != nil {
		base.RedFatal(err.Error())
	}

	ticker := time.NewTicker(time.Second)

	for range ticker.C {
		base.Green("work with tId:%d", tId)
	}
}
