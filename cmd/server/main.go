package main

import (
	"fmt"

	"github.com/rusinov-artem/metrics/server"
	"github.com/spf13/cobra"
)

func NewServerCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "Best metrics project ever",
		Short: "Run server on port 8080",
		Long:  "Run metrics collector server on port 8080",
		Run: func(*cobra.Command, []string) {
			handler := server.NewHandler()
			server.New(handler).Run()
		},
	}

	return rootCmd
}

func main() {
	err := NewServerCmd().Execute()
	if err != nil {
		fmt.Println(err)
	}
}
