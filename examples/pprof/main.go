package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
		if err != http.ErrServerClosed {
			panic(err)
		}
	}()

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		err := server.Shutdown(ctx)
		if err != nil {
			fmt.Printf("shutdown server error:%v", err)
		}
	}()

	var sig = make(chan os.Signal, 1)
	signal.Notify(sig)

	for {
		val := <-sig
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
			signal.Stop(sig)
			return
		}
	}
}
