package config

import (
	"log"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

type Config struct {
	Address         string
	StoreInterval   int
	FileStoragePath string
	Restore         bool
	RestoreString   string
}

func New(cmd *cobra.Command) *Config {
	return fromEnv().FromCli(cmd)
}

func (c *Config) FromCli(cmd *cobra.Command) *Config {
	if c.Address == "" {
		cmd.Flags().StringVarP(
			&c.Address,
			"address",
			"a",
			"0.0.0.0:8080",
			"set address for server to listen on",
		)
	}

	if c.StoreInterval == 0 {
		cmd.Flags().IntVarP(
			&c.StoreInterval,
			"store interval",
			"i",
			300,
			"set store interval",
		)
	}

	if c.FileStoragePath == "" {
		cmd.Flags().StringVarP(
			&c.FileStoragePath,
			"file storage path",
			"f",
			"/tmp/metrics-db.json",
			"set file storage path",
		)
	}

	if c.RestoreString == "" {
		cmd.Flags().BoolVarP(
			&c.Restore,
			"restore on start",
			"r",
			true,
			"enable restore on server start",
		)
	}

	c.Restore, _ = stringToBool(c.RestoreString)

	return c
}

func fromEnv() *Config {
	cfg := &Config{
		Address: func() string {
			addr := os.Getenv("ADDRESS")
			if addr != "" {
				log.Println("Got ADDRESS env variable")
			}
			return addr
		}(),

		StoreInterval: func() int {
			interval := os.Getenv("STORE_INTERVAL")
			val := 0
			if interval != "" {
				val, _ = strconv.Atoi(interval)
				log.Println("Got STORE_INTERVAL env variable")
			}
			return val
		}(),

		FileStoragePath: func() string {
			path := os.Getenv("FILE_STORAGE_PATH")
			if path != "" {
				log.Println("Got FILE_STORAGE_PATH env variable")
			}
			return path
		}(),

		RestoreString: func() string {
			e := os.Getenv("RESTORE")
			if e == "" {
				return ""
			}

			_, s := stringToBool(e)
			return s

		}(),
	}

	if cfg.RestoreString != "" {
		cfg.Restore = false
		if cfg.RestoreString == "true" {
			cfg.Restore = true
		}
	}

	return cfg
}

func stringToBool(e string) (bool, string) {
	e = strings.ToLower(e)
	trueStringList := []string{"true", "1", "t", "on"}
	if slices.Contains(trueStringList, e) {
		return true, "true"
	}

	return false, "false"
}
