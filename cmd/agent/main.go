package main

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/rusinov-artem/metrics/agent"
	"github.com/rusinov-artem/metrics/agent/client"
	"github.com/rusinov-artem/metrics/cmd/agent/config"

	"net/http"
	_ "net/http/pprof"
)

var runAgent = func(cfg *config.Config) {
	ctx := context.Background()
	client := client.New(fmt.Sprintf("http://%s", cfg.Address))
	client.Key = cfg.Key

	logger, _ := zap.NewDevelopment()
	logger.Info("config", zap.Any("config", cfg))

	agent.New(
		client,
		time.Second*time.Duration(cfg.PollInterval),
		time.Second*time.Duration(cfg.ReportInterval),
		cfg.RateLimit,
	).Run(ctx)
}

func NewAgent() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "Agent of best metrics project ever",
		Short: "Run agent to send metrics",
		Long:  "Run agent to send metrics",
	}

	cfg := config.New(cmd)

	cmd.Run = func(*cobra.Command, []string) {
		runAgent(cfg)
	}

	return cmd
}

func main() {
	go http.ListenAndServe(":9999", nil)
	err := NewAgent().Execute()
	if err != nil {
		fmt.Println(err)
	}
}
