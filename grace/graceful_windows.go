package grace

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/grpc-boot/base/v2/logger"
	"github.com/grpc-boot/base/v2/utils"
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
		logger.ZapInfo("receive signal",
			logger.Pid(pid),
			logger.Signal(sig.(syscall.Signal)),
		)

		switch sig {
		case pprofSig:
			if utils.PprofIsRun() {
				logger.ZapInfo("stop pprof",
					logger.Pid(pid),
					logger.Addr(g.pprofAddr),
				)

				if err = utils.StopPprofWithTimeout(5); err != nil {
					logger.ZapError("stop pprof failed",
						logger.Error(err),
						logger.Pid(pid),
						logger.Addr(g.pprofAddr),
					)
				}
			} else if g.pprofAddr != "" {
				go func() {
					logger.ZapInfo("start pprof",
						logger.Pid(pid),
						logger.Addr(g.pprofAddr),
					)

					if err = utils.StartPprof(g.pprofAddr, nil); err != nil {
						if err == http.ErrServerClosed {
							logger.ZapError("pprof Server closed",
								logger.Pid(pid),
								logger.Addr(g.pprofAddr),
							)
							return
						}

						logger.ZapError("start pprof failed",
							logger.Error(err),
							logger.Pid(pid),
							logger.Addr(g.pprofAddr),
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

	logger.ZapInfo("shutdown Server",
		logger.Pid(pid),
		logger.Addr(g.addr),
	)

	if err = g.shutdownTimeout(time.Second * 10); err != nil {
		logger.ZapError("shutdown timeout",
			logger.Error(err),
			logger.Pid(pid),
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
				logger.ZapError("on start error",
					logger.Error(err),
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
