//go:build !windows

package grace

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/grpc-boot/base/v3/logger"
	"github.com/grpc-boot/base/v3/utils"

	"go.uber.org/zap"
)

var (
	EnvKey = `BASE_GRACEFUL`
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
	if os.Getenv(EnvKey) != "" {
		g.isChild = true
		f := os.NewFile(uintptr(3), "")
		l, err := net.FileListener(f)
		if err != nil {
			return err
		}

		g.addr = addr
		g.listener = l
		return nil
	}

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	g.addr = addr
	g.listener = l
	return nil
}

func (g *Graceful) fork() (err error) {
	var (
		args []string
		env  = append(
			os.Environ(),
			fmt.Sprintf("%s=on", EnvKey),
		)
		files = make([]*os.File, 1)
		path  = os.Args[0]
	)

	if len(os.Args) > 1 {
		args = os.Args[1:]
	}

	if len(g.envs) > 0 {
		env = append(env, g.envs...)
	}

	listenerFile, err := g.listener.(*net.TCPListener).File()
	if err != nil {
		return err
	}

	files[0] = listenerFile

	cmd := exec.Command(path, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.ExtraFiles = files
	cmd.Env = env

	err = cmd.Start()
	if err != nil {
		return err
	}

	go func() {
		if err = cmd.Wait(); err != nil {
			logger.Error("child process exit",
				zap.NamedError("Error", err),
				zap.Int("Pid", cmd.Process.Pid),
				zap.String("Cmd", cmd.String()),
			)
		}
	}()

	return
}

func (g *Graceful) initSignals() {
	var (
		err     error
		sig     os.Signal
		sigChan = make(chan os.Signal, 1)
		pid     = syscall.Getpid()
		ok      = true
	)

	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGTERM, syscall.SIGINT, syscall.SIGTSTP)

	for ok {
		sig, ok = <-sigChan
		logger.Info("receive signal",
			zap.Int("Pid", pid),
			zap.Int("Signal", int(sig.(syscall.Signal))),
		)

		switch sig {
		case syscall.SIGHUP:
			if err = g.fork(); err != nil {
				logger.Error("fork failed",
					zap.NamedError("Error", err),
					zap.Int("Pid", pid),
				)
			}
		case syscall.SIGUSR1:
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
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGTSTP:
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

	if g.isChild {
		if err := syscall.Kill(syscall.Getppid(), syscall.SIGTERM); err != nil {
			logger.Error("kill parent process failed",
				zap.NamedError("Error", err),
			)
		}
	}

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
