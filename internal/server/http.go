package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type Server struct {
	host    string
	port    int
	server  *http.Server
	handler *Handler
	service EnrichService
	log     *logrus.Entry
}

func New(host string, port int, service EnrichService, log *logrus.Logger) *Server {
	h := NewHandler(service, log)

	http.HandleFunc("/enrich", h.enrich)

	s := Server{
		host:    host,
		port:    port,
		handler: h,
		service: service,
		log:     log.WithField("module", "http"),
	}

	s.server = &http.Server{
		Addr:              fmt.Sprintf("%s:%d", host, port),
		ReadHeaderTimeout: 30 * time.Second,
	}

	return &s
}

func (s *Server) Run(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	defer s.log.Info("Server is stopped")

	go func() {
		<-ctx.Done()

		err := s.server.Shutdown(shutdownCtx)
		if err != nil {
			s.log.Warningf("server.Shutdown: %s", err)
		}
	}()

	s.log.Infof("Server is running at port %d", s.port)

	err := s.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("ListenAndServe: %w", err)
	}

	return nil
}
