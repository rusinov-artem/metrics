package config

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

type Config struct {
	Address string
}

func (c *Config) FromCli(cmd *cobra.Command) *Config {
	if c.Address == "" {
		cmd.Flags().StringVarP(&c.Address, "address", "a", "localhost:8080", "set address for server to listen on")
	}
	return c
}

func FromEnv() *Config {
	return &Config{
		Address: func() string {
			addr := os.Getenv("ADDRESS")
			if addr != "" {
				log.Println("Got ADDRESS env variable")
			}
			return addr
		}(),
	}
}
