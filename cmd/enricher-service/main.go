package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/AlexZav1327/name-enricher/internal/age"
	"github.com/AlexZav1327/name-enricher/internal/country"
	"github.com/AlexZav1327/name-enricher/internal/gender"
	"github.com/AlexZav1327/name-enricher/internal/server"
	"github.com/AlexZav1327/name-enricher/internal/service"
	"github.com/AlexZav1327/name-enricher/internal/storage"
	_ "github.com/jackc/pgx/v5/stdlib"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/sirupsen/logrus"
)

const (
	host = ""
	port = 8086
)

func main() {
	pgDSN := getEnv(os.Getenv("PG_DSN"), "postgres://user:secret@localhost:5436/postgres?sslmode=disable")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	logger := logrus.StandardLogger()

	pg, err := storage.ConnectDB(ctx, pgDSN, logger)
	if err != nil {
		logger.Panicf("ConnectDB: %s", err)
	}

	err = pg.Migrate(migrate.Up)
	if err != nil {
		logger.Panicf("Migrate: %s", err)
	}

	ageEnrich := age.New(logger)
	genderEnrich := gender.New(logger)
	countryEnrich := country.New(logger)
	enricherService := service.New(pg, ageEnrich, genderEnrich, countryEnrich, logger)
	s := server.New(host, port, enricherService, logger)

	err = s.Run(ctx)
	if err != nil {
		logger.Panicf("Run: %s", err)
	}
}

func getEnv(env, defaultValue string) string {
	value := os.Getenv(env)
	if value == "" {
		return defaultValue
	}

	return value
}
