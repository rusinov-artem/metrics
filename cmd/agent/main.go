package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
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

	cfg := config{
		address: func() string {
			addr := os.Getenv("ADDRESS")
			if addr != "" {
				log.Println("Got ADDRESS env variable")
			}
			return addr
		}(),

		pollInterval: func() int {
			v, _ := strconv.Atoi(os.Getenv("POLL_INTERVAL"))
			if v > 0 {
				log.Println("Got POLL_INTERVAL env variable")
			}
			return v
		}(),

		reportInterval: func() int {
			v, _ := strconv.Atoi(os.Getenv("REPORT_INTERVAL"))
			if v > 0 {
				log.Println("Got REPORT_INTERVAL env variable")
			}
			return v
		}(),
	}

	if cfg.address == "" {
		rootCmd.Flags().StringVarP(&cfg.address, "address", "a", "localhost:8080", "server addres to send metrics to")
	}

	if cfg.pollInterval == 0 {
		rootCmd.Flags().IntVarP(&cfg.pollInterval, "poll_interval", "p", 2, "poll interval")
	}

	if cfg.reportInterval == 0 {
		rootCmd.Flags().IntVarP(&cfg.reportInterval, "report_interval", "r", 10, "report interval")
	}

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
