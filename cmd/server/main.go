package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"

	"github.com/rusinov-artem/metrics/cmd/server/config"
	"github.com/rusinov-artem/metrics/server"
	"github.com/rusinov-artem/metrics/server/handler"
	"github.com/rusinov-artem/metrics/server/middleware"
	"github.com/rusinov-artem/metrics/server/router"
	"github.com/rusinov-artem/metrics/server/storage"

	"github.com/spf13/cobra"
)

var runServer = func(cfg *config.Config) {
	logger, _ := zap.NewDevelopment()
	router := router.New()
	storage, destructor := storage.NewBufferedFileStorage(
		logger,
		cfg.FileStoragePath,
		cfg.Restore,
		cfg.StoreInterval,
	)
	defer destructor()

	dbpool, err := pgxpool.New(context.Background(), cfg.DatabaseDsn)
	if err != nil {
		logger.Error("unable to connect to database", zap.Error(err))
	} else {
		migrate(logger, dbpool)
	}
	defer dbpool.Close()

	handler.New(storage, dbpool).RegisterIn(router)
	router.AddMiddleware(middleware.Logger(logger))
	router.AddMiddleware(middleware.GzipEncoder())
	server.New(router.Mux(), cfg.Address).Run()

}

func NewServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "Best metrics project ever",
		Short: "Run server on port 8080",
		Long:  "Run metrics collector server on port 8080",
	}

	cfg := config.New(cmd)

	cmd.Run = func(*cobra.Command, []string) {
		runServer(cfg)
	}

	return cmd
}

func main() {
	err := NewServerCmd().Execute()
	if err != nil {
		fmt.Println(err)
	}
}

func migrate(log *zap.Logger, pool *pgxpool.Pool) {
	db := stdlib.OpenDBFromPool(pool)
	if err := goose.SetDialect("pgx"); err != nil {
		log.Error("goose: unable to set dialect", zap.Error(err))
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Error("goose: unable get current working directory", zap.Error(err))
	}

	log.Info("goose: migrations dir", zap.String("dir", dir))

	err = goose.Up(db, "./server/migration")
	if err != nil {
		log.Error("goose: unable to run migration", zap.Error(err))
	}
}
