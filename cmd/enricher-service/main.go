package main

import (
	"context"
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
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")

	if err := viper.BindEnv("database.dsn", "PG_DSN"); err != nil {
		logrus.Warningf("viper.BindEnv(): %s", err)
	}

	if err := viper.ReadInConfig(); err != nil {
		logrus.Panicf("viper.ReadInConfig(): %s", err)
	}

	var (
		pgDSN = viper.GetString("database.dsn")
		host  = viper.GetString("server.host")
		port  = viper.GetInt("server.port")
	)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	logger := logrus.StandardLogger()

	pg, err := storage.ConnectDB(ctx, pgDSN, logger)
	if err != nil {
		logger.Panicf("storage.ConnectDB(ctx, pgDSN, logger): %s", err)
	}

	if err = pg.Migrate(migrate.Up); err != nil {
		logger.Panicf("pg.Migrate(migrate.Up): %s", err)
	}

	ageEnrich := age.New(logger)
	genderEnrich := gender.New(logger)
	countryEnrich := country.New(logger)
	enricherService := service.New(pg, ageEnrich, genderEnrich, countryEnrich, logger)
	s := server.New(host, port, enricherService, logger)

	if err = s.Run(ctx); err != nil {
		logger.Panicf("s.Run(ctx): %s", err)
	}
}
