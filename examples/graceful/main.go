package main

import (
	"flag"
	"math/rand"
	"net/http"
	"time"

	"github.com/grpc-boot/base/v2/grace"
	"github.com/grpc-boot/base/v2/logger"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

var useFasthttp int

func main() {
	flag.IntVar(&useFasthttp, "f", 1, "-f")
	flag.Parse()
	var (
		s grace.Serve
	)

	if useFasthttp == 1 {
		logger.Info("init with fasthttp server")

		s = fasthttpServer()
	} else {
		logger.Info("init with http server")

		s = &grace.Server{
			httpServer(),
		}
	}

	g := grace.New(s, nil)
	err := g.Serve(":8080", ":8090")
	if err != nil {
		logger.Error("serve failed",
			zap.NamedError("Error", err),
		)
		panic(err)
	}
}

func fasthttpServer() *fasthttp.Server {
	return &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			duration := time.Duration(rand.Int63n(100)) * time.Millisecond
			time.Sleep(duration)
			_, _ = ctx.Write([]byte("Hello World"))
		},
	}
}

func httpServer() *http.Server {
	return &http.Server{
		Handler: mux(),
	}
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
