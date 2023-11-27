package components

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-boot/base/v2/utils"
)

var (
	ErrNotTcpListener = errors.New("listener is not tcp listener")
)

type GracefulHttp struct {
	server   *http.Server
	listener net.Listener
}

func NewGracefulHttp(server *http.Server) *GracefulHttp {
	return &GracefulHttp{
		server: server,
	}
}

func (hs *GracefulHttp) ShutdownWithTimeout(timeoutSecond int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeoutSecond))
	defer cancel()

	return hs.Shutdown(ctx)
}

func (hs *GracefulHttp) Shutdown(ctx context.Context) error {
	return hs.server.Shutdown(ctx)
}

func (hs *GracefulHttp) Listen() (err error) {
	graceful := flag.Bool("graceful", false, "listen on fd open 3 (internal use only)")
	flag.Parse()

	if *graceful {
		f := os.NewFile(3, "")
		hs.listener, err = net.FileListener(f)
	} else {
		hs.listener, err = net.Listen("tcp", hs.server.Addr)
	}

	if err != nil {
		return err
	}

	go func() {
		if err = hs.server.Serve(hs.listener); err != nil && err != http.ErrServerClosed {
			err = errors.New(fmt.Sprintf("listen error:%v\n", err))
			utils.RedFatal(err.Error())
		}
	}()

	return nil
}

func (hs *GracefulHttp) HandlerSig(pprofAddr string, timeoutSecond int64) {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan)

	for {
		switch <-sigChan {
		case syscall.SIGUSR1:
			subPid, err := hs.Reload()
			if err != nil {
				utils.Red("graceful restart failed with error: %s", err.Error())
				continue
			}

			signal.Stop(sigChan)

			utils.Yellow("graceful restart with sub pid: %d", subPid)
			if err = hs.ShutdownWithTimeout(timeoutSecond); err != nil {
				utils.Red("shutdown with error: %s", err.Error())
			}
			return
		case syscall.SIGUSR2:
			if utils.PprofIsRun() {
				if err := utils.StopPprofWithTimeout(timeoutSecond); err != nil && err != http.ErrServerClosed {
					utils.Red("stop pprof error: %s", err.Error())
				} else {
					utils.Green("stop pprof success")
				}

				continue
			}

			go func() {
				utils.Yellow("start pprof with addr: %s", pprofAddr)
				if err := utils.StartPprof(pprofAddr, nil); err != nil && err != http.ErrServerClosed {
					utils.Red("start pprof error: %s", err.Error())
				}
			}()
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT:
			signal.Stop(sigChan)

			if err := hs.ShutdownWithTimeout(timeoutSecond); err != nil && err != http.ErrServerClosed {
				utils.Red("shutdown with error: %s", err.Error())
			} else {
				utils.Green("shutdown success")
			}

			return
		}
	}
}

func (hs *GracefulHttp) Reload() (pid int, err error) {
	tl, ok := hs.listener.(*net.TCPListener)
	if !ok {
		return 0, ErrNotTcpListener
	}

	f, err := tl.File()
	if err != nil {
		return 0, err
	}

	args := []string{"-graceful"}
	cmd := exec.Command(os.Args[0], args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.ExtraFiles = []*os.File{f}
	if err = cmd.Start(); err != nil {
		return 0, err
	}

	return cmd.Process.Pid, nil
}
