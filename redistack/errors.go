package redistack

import (
	"strings"

	"github.com/redis/rueidis"
)

func ErrBusyGroup(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "BUSYGROUP")
}

func ErrNil(err error) bool {
	return err == rueidis.Nil
}

func IsNil(err error) bool {
	return err == nil || ErrNil(err)
}
