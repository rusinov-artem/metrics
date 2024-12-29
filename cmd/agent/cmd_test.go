package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rusinov-artem/metrics/cmd/agent/config"
)

func Test_CanHandleCommandLineArgs(t *testing.T) {
	runAgent = func(cfg *config.Config) {
		assert.Equal(t, "test_address", cfg.Address)
		assert.Equal(t, 60, cfg.PollInterval)
		assert.Equal(t, 90, cfg.ReportInterval)
		assert.Equal(t, 99, cfg.RateLimit)
		assert.Equal(t, "some_key", cfg.Key)
		assert.Equal(t, "public.pem", cfg.CryptoKey)
	}
	cmd := NewAgent()
	cmd.SetArgs([]string{"",
		"-a", "test_address",
		"-p", "60",
		"-r", "90",
		"-k", "some_key",
		"-l", "99",
		"--crypto-key", "public.pem",
	})
	err := cmd.Execute()
	assert.NoError(t, err)
}

func Test_CanGetValuesFromEnv(t *testing.T) {
	_ = os.Setenv("ADDRESS", "test_address")
	_ = os.Setenv("POLL_INTERVAL", "60")
	_ = os.Setenv("REPORT_INTERVAL", "90")
	_ = os.Setenv("KEY", "some_key_from_env")
	_ = os.Setenv("RATE_LIMIT", "100")
	_ = os.Setenv("CRYPTO_KEY", "public.pem")
	runAgent = func(cfg *config.Config) {
		assert.Equal(t, "test_address", cfg.Address)
		assert.Equal(t, 60, cfg.PollInterval)
		assert.Equal(t, 90, cfg.ReportInterval)
		assert.Equal(t, "some_key_from_env", cfg.Key)
		assert.Equal(t, 100, cfg.RateLimit)
		assert.Equal(t, "public.pem", cfg.CryptoKey)
	}
	cmd := NewAgent()
	err := cmd.Execute()
	assert.NoError(t, err)
}
