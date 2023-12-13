package utils

import (
	"fmt"
)

func FiveMinuteAxis() []string {
	var (
		hour   = 0
		minute = 0
		index  = 0
		axis   = make([]string, 288)
	)

	for hour < 24 {
		for minute < 60 {
			axis[index] = fmt.Sprintf("%02d:%02d", hour, minute)
			index++
			minute += 5
		}

		hour++
		minute = 0
	}

	return axis
}

func MinuteAxis() []string {
	var (
		hour   = 0
		minute = 0
		index  = 0
		axis   = make([]string, 1440)
	)

	for hour < 24 {
		for minute < 60 {
			axis[index] = fmt.Sprintf("%02d:%02d", hour, minute)
			index++
			minute++
		}

		hour++
		minute = 0
	}

	return axis
}
