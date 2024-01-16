package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type Server struct {
	host    string
	port    int
	server  *http.Server
	handler *Handler
	service EnricherService
	log     *logrus.Entry
}

func New(host string, port int, service EnricherService, log *logrus.Logger) *Server {
	h := NewHandler(service, log)

	s := Server{
		host:    host,
		port:    port,
		handler: h,
		service: service,
		log:     log.WithField("module", "http"),
	}

	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Route("/api/v1", func(r chi.Router) {
			r.Post("/user/enrich", h.enrich)
		})
	})

	s.server = &http.Server{
		Addr:              fmt.Sprintf("%s:%d", host, port),
		Handler:           r,
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
