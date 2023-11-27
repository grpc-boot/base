package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/grpc-boot/base/v2/components"
	"github.com/grpc-boot/base/v2/utils"
)

func main() {
	var (
		server = &http.Server{
			Addr:    ":8080",
			Handler: mux(),
		}

		gracehttp = components.NewGracefulHttp(server)
	)

	if err := gracehttp.Listen(); err != nil {
		utils.RedFatal("start server failed with error:%s", err.Error())
	}

	gracehttp.HandlerSig(":8090", 10)
}

func mux() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		duration := time.Duration(rand.Int63n(100)) * time.Millisecond
		time.Sleep(duration)
		_, _ = w.Write([]byte("Hello World"))
	})

	return router
}
