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
	DatabaseDsn     string
	Restore         bool
	RestoreString   string
	Key             string
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
	} else {
		cmd.Flags().BoolVarP(
			&c.Restore,
			"restore on start",
			"r",
			false,
			"enable restore on server start",
		)
	}

	if c.DatabaseDsn == "" {
		cmd.Flags().StringVarP(
			&c.DatabaseDsn,
			"database dsn",
			"d",
			"",
			"enable restore on server start",
		)
	}

	cmd.Flags().StringVarP(
		&c.Key,
		"hash key",
		"k",
		os.Getenv("KEY"),
		"hash key to check request sign",
	)

	c.Restore, _ = stringToBool(c.RestoreString)

	return c
}

func fromEnv() *Config {
	cfg := &Config{
		Address: func() string {
			v := os.Getenv("ADDRESS")
			if v != "" {
				log.Println("Got ADDRESS env variable")
			}
			return v
		}(),

		StoreInterval: func() int {
			v := os.Getenv("STORE_INTERVAL")
			val := 0
			if v != "" {
				val, _ = strconv.Atoi(v)
				log.Println("Got STORE_INTERVAL env variable")
			}
			return val
		}(),

		FileStoragePath: func() string {
			v := os.Getenv("FILE_STORAGE_PATH")
			if v != "" {
				log.Println("Got FILE_STORAGE_PATH env variable")
			}
			return v
		}(),

		RestoreString: func() string {
			v := os.Getenv("RESTORE")
			if v == "" {
				return ""
			}

			_, s := stringToBool(v)
			return s

		}(),

		DatabaseDsn: func() string {
			v := os.Getenv("DATABASE_DSN")
			if v != "" {
				log.Println("Got DATABASE_DSN env variable")
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
