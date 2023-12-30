package httpserver

import (
	"context"
	"net/http"
	"this_module/config"
	"time"
)

type Server struct {
	Server          *http.Server
	Notify          chan error
	ShutdownTimeout time.Duration
}

func (s *Server) Start() {
	s.Notify <- s.Server.ListenAndServe()
	close(s.Notify)
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), s.ShutdownTimeout*time.Second)
	defer cancel()

	s.Server.Shutdown(ctx)
}

func New(router http.Handler, cfg *config.Config) *Server {
	httpServer := &http.Server{
		Addr:    cfg.HTTP.Port,
		Handler: router,
	}

	s := &Server{
		Server:          httpServer,
		Notify:          make(chan error, 1),
		ShutdownTimeout: time.Duration(cfg.WaitingClose),
	}

	go s.Start()

	return s
}
