package main

import (
	"context"
	"fmt"
	"time"

	"github.com/rusinov-artem/metrics/agent"
	"github.com/rusinov-artem/metrics/agent/client"
	"github.com/spf13/cobra"
)

func NewAgent() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "Best metrics project ever",
		Short: "Run agent to send metrics",
		Long:  "Run agent to send metrics",
		Run: func(*cobra.Command, []string) {
			ctx := context.Background()
			client := client.New("http://localhost:8080")
			agent.New(client, time.Second*2, time.Second*10).Run(ctx)
		},
	}

	return rootCmd
}

func main() {
	err := NewAgent().Execute()
	if err != nil {
		fmt.Println(err)
	}
}
