package grace

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/grpc-boot/base/v3/logger"
	"github.com/grpc-boot/base/v3/utils"

	"go.uber.org/zap"
)

type Graceful struct {
	isChild       bool
	server        Serve
	addr          string
	pprofAddr     string
	listener      net.Listener
	envs          []string
	onShutdown    Handler
	onStart       Handler
	signalHandler func(sig os.Signal)
}

func NewWithHttpServer(server *http.Server, envs []string) *Graceful {
	return New(&Server{server}, envs)
}

func New(server Serve, envs []string) *Graceful {
	g := Graceful{
		server: server,
	}

	g.envs = envs
	return &g
}

func (g *Graceful) initListener(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	g.addr = addr
	g.listener = l
	return nil
}

func (g *Graceful) initSignals() {
	var (
		err     error
		sig     os.Signal
		sigChan = make(chan os.Signal, 1)
		pid     = syscall.Getpid()
		ok      = true
	)

	pprofSig := syscall.Signal(10)
	signal.Notify(sigChan, syscall.SIGHUP, pprofSig, syscall.SIGTERM, syscall.SIGINT)

	for ok {
		sig, ok = <-sigChan
		logger.Info("receive signal",
			zap.Int("Pid", pid),
			zap.Int("Signal", int(sig.(syscall.Signal))),
		)

		switch sig {
		case pprofSig:
			if utils.PprofIsRun() {
				logger.Info("stop pprof",
					zap.Int("Pid", pid),
					zap.String("Addr", g.pprofAddr),
				)

				if err = utils.StopPprofWithTimeout(5); err != nil {
					logger.Error("stop pprof failed",
						zap.NamedError("Error", err),
						zap.Int("Pid", pid),
						zap.String("Addr", g.pprofAddr),
					)
				}
			} else if g.pprofAddr != "" {
				go func() {
					logger.Info("start pprof",
						zap.Int("Pid", pid),
						zap.String("Addr", g.pprofAddr),
					)

					if err = utils.StartPprof(g.pprofAddr, nil); err != nil {
						if errors.Is(err, http.ErrServerClosed) {
							logger.Error("pprof Server closed",
								zap.Int("Pid", pid),
								zap.String("Addr", g.pprofAddr),
							)
							return
						}

						logger.Error("start pprof failed",
							zap.NamedError("Error", err),
							zap.Int("Pid", pid),
							zap.String("Addr", g.pprofAddr),
						)
					}
				}()
			}
		case syscall.SIGINT, syscall.SIGTERM:
			signal.Stop(sigChan)
			close(sigChan)
			ok = false
			break
		default:
			if g.signalHandler != nil {
				g.signalHandler(sig)
			}
		}
	}

	logger.Info("shutdown Server",
		zap.Int("Pid", pid),
		zap.String("Addr", g.pprofAddr),
	)

	if err = g.shutdownTimeout(time.Second * 10); err != nil {
		logger.Error("shutdown timeout",
			zap.NamedError("Error", err),
			zap.Int("Pid", pid),
			zap.String("Addr", g.pprofAddr),
		)
	}
}

func (g *Graceful) OnStart(handler Handler) {
	g.onStart = handler
}

func (g *Graceful) OnShutdown(handler Handler) {
	g.onShutdown = handler
}

func (g *Graceful) init(pprofAddr string) {
	g.pprofAddr = pprofAddr

	if g.onStart != nil {
		go func() {
			if err := g.onStart(); err != nil {
				logger.Error("on start error",
					zap.NamedError("Error", err),
				)
			}
		}()
	}

	go g.initSignals()
}

func (g *Graceful) Serve(addr, pprofAddr string) (err error) {
	err = g.initListener(addr)
	if err != nil {
		return err
	}

	g.init(pprofAddr)

	return g.server.Serve(g.listener)
}

func (g *Graceful) ServeTLS(addr, pprofAddr string, certFile, keyFile string) (err error) {
	err = g.initListener(addr)
	if err != nil {
		return err
	}

	g.init(pprofAddr)

	return g.server.ServeTLS(g.listener, certFile, keyFile)
}

func (g *Graceful) shutdownTimeout(timeout time.Duration) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)

	if g.onShutdown != nil {
		wg.Add(1)
		go func() {
			err = utils.Timeout(timeout, func(args ...any) {
				err = g.onShutdown()
			})
			wg.Done()
		}()
	}

	go func() {
		err = g.server.ShutdownWithContext(ctx)
		wg.Done()
	}()

	wg.Wait()
	return
}
