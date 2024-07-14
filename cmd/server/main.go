package main

import (
	"fmt"

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
	handler.New(storage).RegisterIn(router)
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
