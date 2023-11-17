package components

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"time"

	"go.uber.org/atomic"
)

var (
	pprofStatus atomic.Bool
	pprofServer *http.Server
)

// PprofIsRun _
func PprofIsRun() bool {
	return pprofStatus.Load()
}

// StartPprof _
func StartPprof(addr string, handler http.Handler) error {
	if !pprofStatus.CAS(false, true) {
		return nil
	}

	pprofServer = &http.Server{
		Handler: handler,
		Addr:    addr,
	}

	err := pprofServer.ListenAndServe()
	if err != nil {
		pprofStatus.Store(false)
	}
	return err
}

// StopPprof _
func StopPprof(ctx context.Context) error {
	if !pprofStatus.CAS(true, false) {
		return nil
	}

	return pprofServer.Shutdown(ctx)
}

// StopPprofWithTimeout _
func StopPprofWithTimeout(seconds int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(seconds))
	defer cancel()

	return StopPprof(ctx)
}
