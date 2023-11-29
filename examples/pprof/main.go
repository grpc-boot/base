package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-boot/base/v2/internal"
	"github.com/grpc-boot/base/v2/utils"
)

type router struct {
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	_, _ = w.Write(internal.String2Bytes(`ok`))
}

func main() {
	server := http.Server{
		Addr:    ":8080",
		Handler: &router{},
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()

	var sig = make(chan os.Signal, 1)
	signal.Notify(sig)

	for {
		val, ok := <-sig
		if !ok {
			break
		}

		switch val {
		case syscall.SIGUSR1:
			if utils.PprofIsRun() {
				continue
			}

			go func() {
				err := utils.StartPprof(":8081", nil)
				if err != nil {
					fmt.Printf("start pprof error:%v", err)
				}
			}()
		case syscall.SIGUSR2:
			if !utils.PprofIsRun() {
				continue
			}

			err := utils.StopPprofWithTimeout(10)
			if err != nil {
				fmt.Printf("stop pprof error:%v", err)
			}
		default:
			close(sig)
			signal.Stop(sig)
			break
		}
	}
}
