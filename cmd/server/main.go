package main

import (
	"fmt"

	"github.com/rusinov-artem/metrics/cmd/server/config"
	"github.com/rusinov-artem/metrics/server"
	"github.com/rusinov-artem/metrics/server/handler"
	"github.com/rusinov-artem/metrics/server/metrics"
	"github.com/rusinov-artem/metrics/server/router"

	"github.com/spf13/cobra"
)

var runServer = func(cfg *config.Config) {
	router := router.New()
	handler.New(metrics.NewInMemory()).RegisterIn(router)
	server.New(router.Mux(), cfg.Address).Run()
}

func NewServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "Best metrics project ever",
		Short: "Run server on port 8080",
		Long:  "Run metrics collector server on port 8080",
	}

	cfg := config.FromEnv().FromCli(cmd)

	cmd.Run = func(_ *cobra.Command, _ []string) {
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
