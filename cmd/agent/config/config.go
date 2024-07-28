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
	Key            string
}

func New(cmd *cobra.Command) *Config {
	return fromEnv().FromCli(cmd)
}

func (c *Config) FromCli(cmd *cobra.Command) *Config {
	if c.Address == "" {
		cmd.Flags().StringVarP(&c.Address, "address", "a", "localhost:8080", "server address to send metrics to")
	}

	if c.PollInterval == 0 {
		cmd.Flags().IntVarP(&c.PollInterval, "poll_interval", "p", 2, "poll interval")
	}

	if c.ReportInterval == 0 {
		cmd.Flags().IntVarP(&c.ReportInterval, "report_interval", "r", 10, "report interval")
	}

	cmd.Flags().StringVarP(&c.Key, "key", "k", os.Getenv("KEY"), "key to sign request")

	return c
}

func fromEnv() *Config {
	return &Config{
		Address: func() string {
			v := os.Getenv("ADDRESS")
			if v != "" {
				log.Println("Got ADDRESS env variable")
			}
			return v
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

		Key: func() string {
			v := os.Getenv("KEY")
			if v != "" {
				log.Printf("Got KEY = %s\n", v)
			}
			return v
		}(),
	}
}
