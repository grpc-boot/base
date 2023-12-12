package monitor

import (
	"errors"
	"golang.org/x/exp/rand"
	"sync"
	"testing"
	"time"

	"github.com/grpc-boot/base/v2/utils"
)

func TestMonitor_Add(t *testing.T) {
	var (
		dOpt = DefaultOptions()
		m    = NewMonitor(dOpt)
		cnt  = 8
	)

	mp := [][]string{
		{
			"GET /v1/user/login", "登录",
		},
		{
			"POST /v1/user/regis", "注册",
		},
		{
			"GET /v1/user/info", "获取用户信息",
		},
	}

	for _, item := range mp {
		m.Path(item[0], item[1])
	}

	for i := 0; i < cnt; i++ {
		go func() {
			for {
				index := rand.Int()
				switch index % 4 {
				case 0:
					m.Add(GaugeRequestCount, mp[rand.Intn(3)][0], utils.OK, 1)
				case 1:
					m.Add(GaugeRequestLen, mp[rand.Intn(3)][0], utils.OK, uint64(rand.Intn(40960)))
				case 2:
					m.Add(GaugeResponseLen, mp[rand.Intn(3)][0], dOpt.CodeList[rand.Intn(len(dOpt.CodeList))], uint64(rand.Intn(40960)))
				case 4:
					m.Add(GaugeResponseCount, mp[rand.Intn(3)][0], dOpt.CodeList[rand.Intn(len(dOpt.CodeList))], 1)
				}

				go func() {
					defer func() {
						if err := recover(); err != nil {
							m.AddGauge(GaugePanicCount, 1)
						}
					}()

					if rand.Int()%10 == 0 {
						panic(errors.New(time.Now().String()))
					}
				}()
				time.Sleep(time.Microsecond * time.Duration(500+rand.Intn(800)))
			}
		}()
	}

	var wg sync.WaitGroup
	wg.Add(1)
	time.AfterFunc(time.Second*30, func() {
		data, _ := utils.JsonEncode(m.Info())
		t.Logf("data: %s", data)

		wg.Done()
	})

	wg.Wait()
}
