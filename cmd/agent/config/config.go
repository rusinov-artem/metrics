package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

const defaultAddres = "localhost:8080"
const defaultPollInterval = 2
const defaultReportInterval = 10
const defaultCryptoKey = ""

type Config struct {
	Address        string
	PollInterval   int
	ReportInterval int
	Key            string
	RateLimit      int
	CryptoKey      string
	Config         string
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

	return (&Config{}).fromEnv().FromCli(cmd).FromCustomSource(cfg)
}

func (c *Config) FromCustomSource(cfg *Config) *Config {
	if cfg == nil {
		return c
	}

	if c.Address == defaultAddres {
		c.Address = cfg.Address
	}

	if c.PollInterval == defaultPollInterval {
		c.PollInterval = cfg.PollInterval
	}

	if c.ReportInterval == defaultReportInterval {
		c.ReportInterval = cfg.ReportInterval
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
		cmd.Flags().StringVarP(&c.Address, "address", "a", defaultAddres, "server address to send metrics to")
	}

	if c.PollInterval == 0 {
		cmd.Flags().IntVarP(&c.PollInterval, "poll_interval", "p", defaultPollInterval, "poll interval")
	}

	if c.ReportInterval == 0 {
		cmd.Flags().IntVarP(&c.ReportInterval, "report_interval", "r", defaultReportInterval, "report interval")
	}

	cmd.Flags().StringVarP(&c.Key, "key", "k", os.Getenv("KEY"), "key to sign request")

	if c.RateLimit == 0 {
		cmd.Flags().IntVarP(&c.RateLimit, "rate_limit", "l", 10, "rate limit")
	}

	if c.CryptoKey == "" {
		cmd.Flags().
			StringVar(&c.CryptoKey, "crypto-key", os.Getenv("CRYPTO_KEY"), "public key to encrypt data")
	}

	return c
}

func (c *Config) fromEnv() *Config {
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

		RateLimit: func() int {
			v := os.Getenv("RATE_LIMIT")
			if v != "" {
				log.Printf("Got RATE_LIMIT = %s\n", v)
			}
			val, _ := strconv.Atoi(v)
			return val
		}(),
	}
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
		Address        string `json:"address"`
		ReportInterval int    `json:"report_interval"`
		PollInterval   int    `json:"poll_interval"`
		CryptoKey      string `json:"crypto_key"`
	}{
		Address:        "localhost:8080",
		ReportInterval: 10,
		PollInterval:   2,
	}

	err = json.Unmarshal(cfgData, &tmp)
	if err != nil {
		return cfg, err
	}

	cfg.Address = tmp.Address
	cfg.ReportInterval = tmp.ReportInterval
	cfg.PollInterval = tmp.PollInterval
	cfg.CryptoKey = tmp.CryptoKey

	return cfg, nil
}
