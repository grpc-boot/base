package base

import (
	"context"
	"net/http"
	_ "net/http/pprof"

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

	return pprofServer.ListenAndServe()
}

// StopPprof _
func StopPprof(ctx context.Context) error {
	if !pprofStatus.CAS(true, false) {
		return nil
	}

	err := pprofServer.Shutdown(ctx)

	pprofStatus.Store(false)
	return err
}
