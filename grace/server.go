package grace

import (
	"context"
	"net"
	"net/http"
)

type Serve interface {
	Serve(ln net.Listener) error
	ServeTLS(ln net.Listener, certFile, keyFile string) error
	ShutdownWithContext(ctx context.Context) (err error)
}

type Server struct {
	*http.Server
}

func (s *Server) ShutdownWithContext(ctx context.Context) (err error) {
	return s.Shutdown(ctx)
}
