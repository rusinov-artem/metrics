package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

const defaultAddress = "localhost:8080"
const defaultRestore = true
const defaultStoreInterval = 300
const defaultStoreFile = "/tmp/metrics-db.json"
const defaultDatabaseDsn = ""
const defaultCryptoKey = ""

type Config struct {
	Address         string
	StoreInterval   int
	FileStoragePath string
	DatabaseDsn     string
	Restore         bool
	RestoreString   string
	Key             string
	CryptoKey       string
	Config          string
}

func New(cmd *cobra.Command) *Config {
	var err error
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	jsonCfgFile := fs.String("config", os.Getenv("CONFIG"), "json config file")
	cfg := (*Config)(nil)
	if *jsonCfgFile != "" {
		cfg, err = FromJSONFile(*jsonCfgFile)
		if err != nil {
			fmt.Printf("unable to fetch config from %q \n", *jsonCfgFile)
		}
	}

	return fromEnv().FromCli(cmd).FromCustomSource(cfg)
}

func (c *Config) FromCustomSource(cfg *Config) *Config {
	if cfg == nil {
		return c
	}

	if c.Address == defaultAddress {
		c.Address = cfg.Address
	}

	if c.StoreInterval == defaultStoreInterval {
		c.StoreInterval = cfg.StoreInterval
	}

	if c.DatabaseDsn == defaultDatabaseDsn {
		c.DatabaseDsn = cfg.DatabaseDsn
	}

	if c.CryptoKey == defaultCryptoKey {
		c.CryptoKey = cfg.CryptoKey
	}

	return c
}

func (c *Config) FromCli(cmd *cobra.Command) *Config {
	if c.Config == "" {
		cmd.Flags().StringVarP(&c.Config, "config", "c", os.Getenv("CONFIG"), "json config file")
	}

	if c.Address == "" {
		cmd.Flags().StringVarP(
			&c.Address,
			"address",
			"a",
			defaultAddress,
			"set address for server to listen on",
		)
	}

	if c.StoreInterval == 0 {
		cmd.Flags().IntVarP(
			&c.StoreInterval,
			"store interval",
			"i",
			defaultStoreInterval,
			"set store interval",
		)
	}

	if c.FileStoragePath == "" {
		cmd.Flags().StringVarP(
			&c.FileStoragePath,
			"file storage path",
			"f",
			defaultStoreFile,
			"set file storage path",
		)
	}

	if c.RestoreString == "" {
		cmd.Flags().BoolVarP(
			&c.Restore,
			"restore on start",
			"r",
			defaultRestore,
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

	if c.CryptoKey == "" {
		cmd.Flags().StringVar(&c.CryptoKey, "crypto-key", os.Getenv("CRYPTO_KEY"), "private key to decode messages")
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
func FromJSONFile(cfgFile string) (*Config, error) {
	cfg := &Config{}
	f, err := os.Open(cfgFile)
	if err != nil {
		return cfg, err
	}

	cfgData, err := io.ReadAll(f)
	if err != nil {
		return cfg, err
	}

	tmp := struct {
		Address       string `json:"address"`
		Restore       bool   `json:"restore"`
		StoreInterval int    `json:"store_interval"`
		StoreFile     string `json:"store_file"`
		DatabaseDSN   string `json:"database_dsn"`
		CryptoKey     string `json:"crypto_key"`
	}{
		Address:       defaultAddress,
		Restore:       defaultRestore,
		StoreInterval: defaultStoreInterval,
		StoreFile:     defaultStoreFile,
		DatabaseDSN:   defaultDatabaseDsn,
		CryptoKey:     defaultCryptoKey,
	}

	err = json.Unmarshal(cfgData, &tmp)
	if err != nil {
		return cfg, err
	}

	cfg.Address = tmp.Address

	return cfg, nil
}
