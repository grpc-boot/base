package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"golang.org/x/exp/rand"
	"net"
	"time"

	"github.com/grpc-boot/base/v2/components"
	"github.com/grpc-boot/base/v2/gored"
	"github.com/grpc-boot/base/v2/grace"
	"github.com/grpc-boot/base/v2/logger"
	"github.com/grpc-boot/base/v2/monitor"
	"github.com/grpc-boot/base/v2/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap/zapcore"
)

var (
	m *monitor.Monitor
)

func init() {
	rand.Seed(uint64(time.Now().Unix()))

	err := logger.InitZapWithOption(logger.Option{
		Level:      int8(zapcore.InfoLevel),
		Path:       "./",
		TickSecond: 5,
		MaxDays:    1,
	})

	if err != nil {
		panic(err)
	}

	opt := monitor.DefaultOptions()
	m = monitor.NewMonitor(opt)
}

func BindRoute(router fiber.Router, method, path, desc string, handlers ...fiber.Handler) fiber.Router {
	mPath := fmt.Sprintf("%s %s", method, path)
	m.Path(mPath, desc)

	router.Add(method, path, handlers...)

	return router
}

func RetJson(ctx *fiber.Ctx, sts *components.Status) (length int, err error) {
	defer sts.Close()

	length, err = ctx.Write(sts.JsonMarshal())

	var (
		path   = fmt.Sprintf("%s %s", ctx.Method(), ctx.Request().URI().PathOriginal())
		outLen = uint64(len(ctx.Response().Header.Header()) + len(ctx.Response().Body()))
	)

	m.AddWithStatus(monitor.GaugeResponseCount, path, sts, 1)
	m.AddWithStatus(monitor.GaugeResponseLen, path, sts, outLen)

	return
}

func main() {
	var (
		engine = fiber.New()
		opt    = gored.DefaultOptions()
		red    = redis.NewClient(&opt)
	)

	go func() {
		ticker := time.NewTicker(time.Second * 60)
		for range ticker.C {
			var (
				prefix              = time.Now().Format("20060102")
				info                = m.Info()
				gaugeKeys, codeKeys = info.Keys(prefix)
				cmdTimeout          = time.Second * 3
			)

			for _, key := range gaugeKeys {
				gored.TimeoutDo(cmdTimeout, func(ctx context.Context) {
					var (
						cmd = red.HGetAll(ctx, key)
						err = gored.DealCmdErr(cmd)
					)

					if err != nil {
						logger.ZapError("get monitor info failed",
							logger.Error(err),
						)
						return
					}

					logger.ZapDebug("get monitor info",
						logger.Key(key),
						logger.Value(cmd.Val()),
					)
				})
			}

			for _, key := range codeKeys {
				gored.TimeoutDo(cmdTimeout, func(ctx context.Context) {
					var (
						cmd = red.HGetAll(ctx, key)
						err = gored.DealCmdErr(cmd)
					)

					if err != nil {
						logger.ZapError("get monitor info failed",
							logger.Error(err),
						)
						return
					}

					logger.ZapInfo("get monitor info",
						logger.Key(key),
						logger.Value(cmd.Val()),
					)
				})
			}
		}
	}()

	// monitor信息同步到redis
	m.WithStorage(monitor.RedisStorage(red))

	// handler panic
	engine.Use(func(ctx *fiber.Ctx) (err error) {
		defer func() {
			if er := recover(); er != nil {
				m.AddGauge(monitor.GaugePanicCount, 1)
				err = nil
			}
		}()

		return ctx.Next()
	})

	// monitor request
	engine.Use(func(ctx *fiber.Ctx) error {
		var (
			path  = fmt.Sprintf("%s %s", ctx.Method(), ctx.Request().URI().PathOriginal())
			inLen = uint64(len(ctx.Request().Header.Header()) + len(ctx.Request().Body()))
		)

		m.Add(monitor.GaugeRequestCount, path, utils.OK, 1)
		m.Add(monitor.GaugeRequestLen, path, utils.OK, inLen)

		return ctx.Next()
	})

	BindRoute(engine, fiber.MethodGet, "/monitor/info", "监控详情", func(ctx *fiber.Ctx) error {
		_, _ = RetJson(ctx, components.StatusOk(m.Info()))
		return nil
	})

	BindRoute(engine, fiber.MethodGet, "/monitor/axis", "折线数据", func(ctx *fiber.Ctx) error {
		var (
			data     = map[string]any{}
			prefix   = time.Now().Format("20060102")
			axisData = utils.MinuteAxis()
		)
		data["axisData"] = axisData
		data["panicCount"], _ = monitor.GaugeLineFromRedis(red, m, prefix, monitor.GaugePanicCount, axisData)
		data["requestCnt"], data["codeData"], _ = monitor.CodeLineFromRedis(red, m, prefix, monitor.GaugeRequestCount, "GET /user/regis", axisData)

		_, _ = RetJson(ctx, components.StatusOk(data))
		return nil
	})

	BindRoute(engine, fiber.MethodGet, "/user/regis", "注册接口", func(ctx *fiber.Ctx) error {
		if rand.Int()%2 == 0 {
			_, _ = RetJson(ctx, components.StatusOk(map[string]interface{}{
				"id":        rand.Int31(),
				"createdAt": time.Now().Unix(),
			}))
		} else {
			panic(errors.New("recover test at :" + time.Now().String()))
		}

		return nil
	})

	BindRoute(engine, fiber.MethodGet, "/user/login", "登录接口", func(ctx *fiber.Ctx) error {
		_, _ = RetJson(ctx, components.StatusOk(map[string]interface{}{
			"id":      rand.Int31(),
			"loginAt": time.Now().Unix(),
		}))
		return nil
	})

	s := grace.New(&server{engine}, nil)
	if err := s.Serve(":8080", ":8081"); err != nil {
		utils.Red("start server failed with error: %v", err)
		panic(err)
	}
}

type server struct {
	*fiber.App
}

func (s *server) Serve(ln net.Listener) error {
	return s.App.Listener(ln)
}

func (s *server) ServeTLS(ln net.Listener, certFile, keyFile string) error {
	// Check for valid cert/key path
	if len(certFile) == 0 || len(keyFile) == 0 {
		return errors.New("tls: provide a valid cert or key path")
	}

	// Set TLS config with handler
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return fmt.Errorf("tls: cannot load TLS key pair from certFile=%q and keyFile=%q: %w", certFile, keyFile, err)
	}

	tlsHandler := &fiber.TLSHandler{}
	config := &tls.Config{
		MinVersion: tls.VersionTLS12,
		Certificates: []tls.Certificate{
			cert,
		},
		GetCertificate: tlsHandler.GetClientInfo,
	}

	return s.App.Listener(tls.NewListener(ln, config))
}
