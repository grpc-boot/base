package monitor

import (
	"context"
	"strconv"
	"time"

	"github.com/grpc-boot/base/v2/gored"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
)

func GaugeLineFromRedis(red *redis.Client, m *Monitor, prefix, gaugeName string, axisData []string) (seriesData []uint64, err error) {
	seriesData = make([]uint64, len(axisData))
	if len(axisData) < 1 {
		return
	}

	var (
		mInfo = m.Info()
		list  = mInfo.GaugesInfo
	)

	if len(list) < 1 {
		return
	}

	var info Info

	for _, i := range list {
		if gaugeName == i.Path {
			info = i
			break
		}
	}

	if info.Path == "" {
		return
	}

	gored.TimeoutDo(time.Second*6, func(ctx context.Context) {
		cmd := red.HGetAll(ctx, info.Key(mInfo.Name, prefix))
		err = gored.DealCmdErr(cmd)
		if err != nil {
			return
		}

		data := cmd.Val()
		if len(data) < 1 {
			return
		}

		for index, timeStr := range axisData {
			if val, ok := data[timeStr]; ok {
				seriesData[index], _ = strconv.ParseUint(val, 10, 64)
			}
		}
	})

	return
}

func CodeLineFromRedis(red *redis.Client, m *Monitor, prefix, gaugeName, path string, axisData []string) (seriesData []uint64, codeSeriesData map[codes.Code][]uint64, err error) {
	seriesData = make([]uint64, len(axisData))
	if len(axisData) < 1 {
		return
	}

	var (
		mInfo = m.Info()
		mm    = mInfo.CodesInfo
	)

	if len(mm) < 1 {
		return
	}

	list, exists := mm[gaugeName]
	if !exists {
		return
	}

	var info CodeInfo
	for _, i := range list {
		if path == i.Path {
			info = i
			break
		}
	}

	if info.Path == "" {
		return
	}

	codeSeriesData = make(map[codes.Code][]uint64, len(info.Sub))
	for _, sub := range info.Sub {
		codeSeriesData[sub.Code] = make([]uint64, len(axisData))
	}

	gored.TimeoutDo(time.Second*6, func(ctx context.Context) {
		cmd := red.HGetAll(ctx, info.Key(mInfo.Name, prefix))
		err = gored.DealCmdErr(cmd)
		if err != nil {
			return
		}

		data := cmd.Val()
		if len(data) < 1 {
			return
		}

		for index, timeStr := range axisData {
			if val, ok := data[timeStr]; ok {
				seriesData[index], _ = strconv.ParseUint(val, 10, 64)
			}

			for _, sub := range info.Sub {
				if val, ok := data[sub.Field(timeStr)]; ok {
					codeSeriesData[sub.Code][index], _ = strconv.ParseUint(val, 10, 64)
				}
			}
		}
	})

	return
}
