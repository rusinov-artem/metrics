package main

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/rusinov-artem/metrics/agent"
	"github.com/rusinov-artem/metrics/agent/client"
	"github.com/rusinov-artem/metrics/cmd/agent/config"
)

var runAgent = func(cfg *config.Config) {
	ctx := context.Background()
	client := client.New(fmt.Sprintf("http://%s", cfg.Address))
	agent.New(
		client,
		time.Second*time.Duration(cfg.PollInterval),
		time.Second*time.Duration(cfg.ReportInterval),
	).Run(ctx)
}

func NewAgent() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "Best metrics project ever",
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
	err := NewAgent().Execute()
	if err != nil {
		fmt.Println(err)
	}
}
