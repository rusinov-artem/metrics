package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewAgent() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "Best metrics project ever",
		Short: "Run agent to send metrics",
		Long:  "Run agent to send metrics",
		Run: func(*cobra.Command, []string) {
			fmt.Println("Cobra agent alive")
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
