package grace

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/grpc-boot/base/v2/logger"
	"github.com/grpc-boot/base/v2/utils"
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
			logger.ZapError("child process exit",
				logger.Error(err),
				logger.Pid(cmd.Process.Pid),
				logger.Cmd(cmd.String()),
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
		logger.ZapInfo("receive signal",
			logger.Pid(pid),
			logger.Signal(sig.(syscall.Signal)),
		)

		switch sig {
		case syscall.SIGHUP:
			if err = g.fork(); err != nil {
				logger.ZapError("fork failed",
					logger.Error(err),
					logger.Pid(pid),
				)
			}
		case syscall.SIGUSR1:
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

	if g.isChild {
		if err := syscall.Kill(syscall.Getppid(), syscall.SIGTERM); err != nil {
			logger.ZapError("kill parent process failed",
				logger.Error(err),
			)
		}
	}

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
