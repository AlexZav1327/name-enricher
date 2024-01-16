package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/AlexZav1327/name-enricher/internal/repo"
	"github.com/AlexZav1327/name-enricher/internal/server"
	"github.com/AlexZav1327/name-enricher/internal/service"
	"github.com/sirupsen/logrus"
)

const (
	host = ""
	port = 8082
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	logger := logrus.StandardLogger()
	er := repo.New(logger)
	enrichService := service.New(er, logger)
	s := server.New(host, port, enrichService, logger)

	err := s.Run(ctx)
	if err != nil {
		logger.Panicf("Run: %s", err)
	}
}
