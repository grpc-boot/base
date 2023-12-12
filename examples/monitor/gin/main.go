package main

import (
	"errors"
	"fmt"
	"golang.org/x/exp/rand"
	"net/http"
	"time"

	"github.com/grpc-boot/base/v2/components"
	"github.com/grpc-boot/base/v2/components/grace"
	"github.com/grpc-boot/base/v2/monitor"
	"github.com/grpc-boot/base/v2/utils"

	"github.com/gin-gonic/gin"
)

var (
	m *monitor.Monitor
)

func init() {
	rand.Seed(uint64(time.Now().Unix()))

	opt := monitor.DefaultOptions()
	opt.ResetSeconds = 10
	opt.CodeGauges = []string{monitor.GaugeRequestCount, monitor.GaugeResponseCount, monitor.GaugeResponseLen}

	m = monitor.NewMonitor(opt)
}

func BindRoute(router *gin.Engine, method, path, desc string, handlers ...gin.HandlerFunc) *gin.Engine {
	mPath := fmt.Sprintf("%s %s", method, path)
	m.Path(mPath, desc)

	router.Match([]string{method}, path, handlers...)

	return router
}

func RetJson(ctx *gin.Context, sts *components.Status) {
	defer sts.Close()

	ctx.Data(http.StatusOK, gin.MIMEJSON, sts.JsonMarshal())

	var (
		path   = fmt.Sprintf("%s %s", ctx.Request.Method, ctx.FullPath())
		outLen = uint64(ctx.Writer.Size())
	)

	m.AddWithStatus(monitor.GaugeResponseCount, path, sts, 1)
	m.AddWithStatus(monitor.GaugeResponseLen, path, sts, outLen)
}

func main() {
	engine := gin.Default()

	// handler panic
	engine.Use(func(ctx *gin.Context) {
		defer func() {
			if er := recover(); er != nil {
				m.AddGauge(monitor.GaugePanicCount, 1)
			}
		}()

		ctx.Next()
	})

	// monitor request
	engine.Use(func(ctx *gin.Context) {
		var (
			path = fmt.Sprintf("%s %s", ctx.Request.Method, ctx.FullPath())
		)

		m.Add(monitor.GaugeRequestCount, path, utils.OK, 1)

		ctx.Next()
	})

	BindRoute(engine, http.MethodGet, "/monitor/info", "监控详情", func(ctx *gin.Context) {
		RetJson(ctx, components.StatusOk(m.Info()))
		return
	})

	BindRoute(engine, http.MethodGet, "/user/regis", "注册接口", func(ctx *gin.Context) {
		if rand.Int()%2 == 0 {
			RetJson(ctx, components.StatusOk(map[string]interface{}{
				"id":        rand.Int31(),
				"createdAt": time.Now().Unix(),
			}))
		} else {
			panic(errors.New("recover test at :" + time.Now().String()))
		}

		return
	})

	BindRoute(engine, http.MethodGet, "/user/login", "登录接口", func(ctx *gin.Context) {
		RetJson(ctx, components.StatusOk(map[string]interface{}{
			"id":      rand.Int31(),
			"loginAt": time.Now().Unix(),
		}))
		return
	})

	server := &http.Server{
		Handler: engine,
	}

	s := grace.NewWithHttpServer(server, nil)
	if err := s.Serve(":8080", ":8081"); err != nil {
		if err == http.ErrServerClosed {
			return
		}
		utils.Red("start server failed with error: %v", err)
		panic(err)
	}
}
