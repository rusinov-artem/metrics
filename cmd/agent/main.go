package main

import (
	"context"
	"fmt"
	"time"

	"github.com/rusinov-artem/metrics/agent"
	"github.com/rusinov-artem/metrics/agent/client"
	"github.com/spf13/cobra"
)

type config struct {
	address        string
	pollInterval   int
	reportInterval int
}

var runAgent = func(cfg config) {
	ctx := context.Background()
	client := client.New(fmt.Sprintf("http://%s", cfg.address))
	agent.New(
		client,
		time.Second*time.Duration(cfg.pollInterval),
		time.Second*time.Duration(cfg.reportInterval),
	).Run(ctx)
}

func NewAgent() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "Best metrics project ever",
		Short: "Run agent to send metrics",
		Long:  "Run agent to send metrics",
	}

	cfg := config{}

	rootCmd.Flags().StringVarP(&cfg.address, "address", "a", "localhost:8080", "server addres to send metrics to")
	rootCmd.Flags().IntVarP(&cfg.pollInterval, "poll_interval", "p", 2, "poll interval")
	rootCmd.Flags().IntVarP(&cfg.reportInterval, "report_interval", "r", 10, "report interval")

	rootCmd.Run = func(*cobra.Command, []string) {
		runAgent(cfg)
	}

	return rootCmd
}

func main() {
	err := NewAgent().Execute()
	if err != nil {
		fmt.Println(err)
	}
}
