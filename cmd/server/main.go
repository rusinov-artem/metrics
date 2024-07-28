package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/rusinov-artem/metrics/cmd/server/config"
	"github.com/rusinov-artem/metrics/server"
	"github.com/rusinov-artem/metrics/server/handler"
	"github.com/rusinov-artem/metrics/server/middleware"
	"github.com/rusinov-artem/metrics/server/migration"
	"github.com/rusinov-artem/metrics/server/router"
	"github.com/rusinov-artem/metrics/server/storage"

	"github.com/spf13/cobra"
)

var runServer = func(cfg *config.Config) {
	var err error
	var metricsStorage handler.MetricsStorage
	var dbpool *pgxpool.Pool

	logger, _ := zap.NewDevelopment()
	router := router.New()
	metricsStorage, destructor := storage.NewBufferedFileStorage(
		logger,
		cfg.FileStoragePath,
		cfg.Restore,
		cfg.StoreInterval,
	)
	defer destructor()

	if cfg.DatabaseDsn != "" {
		dbpool, err = pgxpool.New(context.Background(), cfg.DatabaseDsn)
		if err != nil {
			logger.Error("unable to connect to database", zap.Error(err))
		} else {
			migration.Migrate(logger, dbpool)
		}
		defer dbpool.Close()

		metricsStorage = storage.NewPgxStorage(dbpool)
	}

	handler.New(logger, metricsStorage, dbpool).RegisterIn(router)
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
