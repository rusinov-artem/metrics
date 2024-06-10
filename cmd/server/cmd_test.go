package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CanHandleCommandLineArgs(t *testing.T) {
	runServer = func(cfg config) {
		assert.Equal(t, "test_address", cfg.address)
	}
	cmd := NewServerCmd()
	cmd.SetArgs([]string{"", "-a", "test_address"})
	err := cmd.Execute()
	assert.NoError(t, err)
}

func Test_CanGetValuesFromEnv(t *testing.T) {
	os.Setenv("ADDRESS", "test_address")
	runServer = func(cfg config) {
		assert.Equal(t, "test_address", cfg.address)
	}
	cmd := NewServerCmd()
	err := cmd.Execute()
	assert.NoError(t, err)
}
