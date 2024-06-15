package config

import (
	"log"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

type Config struct {
	Address        string
	PollInterval   int
	ReportInterval int
}

func (c *Config) FromCli(cmd *cobra.Command) *Config {
	if c.Address == "" {
		cmd.Flags().StringVarP(&c.Address, "Address", "a", "localhost:8080", "server address to send metrics to")
	}

	if c.PollInterval == 0 {
		cmd.Flags().IntVarP(&c.PollInterval, "poll_interval", "p", 2, "poll interval")
	}

	if c.ReportInterval == 0 {
		cmd.Flags().IntVarP(&c.ReportInterval, "report_interval", "r", 10, "report interval")
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

		PollInterval: func() int {
			v, _ := strconv.Atoi(os.Getenv("POLL_INTERVAL"))
			if v > 0 {
				log.Println("Got POLL_INTERVAL env variable")
			}
			return v
		}(),

		ReportInterval: func() int {
			v, _ := strconv.Atoi(os.Getenv("REPORT_INTERVAL"))
			if v > 0 {
				log.Println("Got REPORT_INTERVAL env variable")
			}
			return v
		}(),
	}
}
