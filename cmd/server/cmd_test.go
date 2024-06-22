package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rusinov-artem/metrics/cmd/server/config"
)

func Test_CanHandleCommandLineArgs(t *testing.T) {
	runServer = func(cfg *config.Config) {
		assert.Equal(t, "test_address", cfg.Address)
	}
	cmd := NewServerCmd()
	cmd.SetArgs([]string{"", "-a", "test_address"})
	err := cmd.Execute()
	assert.NoError(t, err)
}

func Test_CanGetValuesFromEnv(t *testing.T) {
	os.Setenv("ADDRESS", "test_address")
	runServer = func(cfg *config.Config) {
		assert.Equal(t, "test_address", cfg.Address)
	}
	cmd := NewServerCmd()
	err := cmd.Execute()
	assert.NoError(t, err)
}
