package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rusinov-artem/metrics/server"
	"github.com/rusinov-artem/metrics/server/handler"
	"github.com/spf13/cobra"
)

type config struct {
	address string
}

var runServer = func(cfg config) {
	handler := handler.NewHandler()
	server.New(handler, cfg.address).Run()
}

func NewServerCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "Best metrics project ever",
		Short: "Run server on port 8080",
		Long:  "Run metrics collector server on port 8080",
	}

	cfg := config{
		address: func() string {
			addr := os.Getenv("ADDRESS")
			if addr != "" {
				log.Println("Got ADDRESS env variable")
			}
			return addr
		}(),
	}

	if cfg.address == "" {
		rootCmd.Flags().StringVarP(&cfg.address, "address", "a", "localhost:8080", "set addres for server to listen on")
	}

	rootCmd.Run = func(_ *cobra.Command, _ []string) {
		runServer(cfg)
	}

	return rootCmd
}

func main() {
	err := NewServerCmd().Execute()
	if err != nil {
		fmt.Println(err)
	}
}
