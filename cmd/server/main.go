package main

import (
	"fmt"

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
		Run: func(cmd *cobra.Command, _ []string) {
			cfg := config{
				address: cmd.Flags().Lookup("address").Value.String(),
			}

			runServer(cfg)
		},
	}

	rootCmd.Flags().StringP("address", "a", "localhost:8080", "set addres for server to listen on")

	return rootCmd
}

func main() {
	err := NewServerCmd().Execute()
	if err != nil {
		fmt.Println(err)
	}
}
