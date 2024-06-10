package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CanHandleCommandLineArgs(t *testing.T) {
	runAgent = func(cfg config) {
		assert.Equal(t, "test_address", cfg.address)
		assert.Equal(t, 60, cfg.pollInterval)
		assert.Equal(t, 90, cfg.reportInterval)
	}
	cmd := NewAgent()
	cmd.SetArgs([]string{"", "-a", "test_address", "-p", "60", "-r", "90"})
	err := cmd.Execute()
	assert.NoError(t, err)
}

func Test_CanGetValuesFromEnv(t *testing.T) {
	os.Setenv("ADDRESS", "test_address")
	os.Setenv("POLL_INTERVAL", "60")
	os.Setenv("REPORT_INTERVAL", "90")
	runAgent = func(cfg config) {
		assert.Equal(t, "test_address", cfg.address)
		assert.Equal(t, 60, cfg.pollInterval)
		assert.Equal(t, 90, cfg.reportInterval)
	}
	cmd := NewAgent()
	err := cmd.Execute()
	assert.NoError(t, err)
}
