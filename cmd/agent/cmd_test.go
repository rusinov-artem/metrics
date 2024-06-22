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
	}
	cmd := NewAgent()
	cmd.SetArgs([]string{"", "-a", "test_address", "-p", "60", "-r", "90"})
	err := cmd.Execute()
	assert.NoError(t, err)
}

func Test_CanGetValuesFromEnv(t *testing.T) {
	_ = os.Setenv("ADDRESS", "test_address")
	_ = os.Setenv("POLL_INTERVAL", "60")
	_ = os.Setenv("REPORT_INTERVAL", "90")
	runAgent = func(cfg *config.Config) {
		assert.Equal(t, "test_address", cfg.Address)
		assert.Equal(t, 60, cfg.PollInterval)
		assert.Equal(t, 90, cfg.ReportInterval)
	}
	cmd := NewAgent()
	err := cmd.Execute()
	assert.NoError(t, err)
}
