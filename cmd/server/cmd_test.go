package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rusinov-artem/metrics/cmd/server/config"
)

func Test_CanHandleCommandLineArgs(t *testing.T) {
	runServer = func(cfg *config.Config) {
		assert.Equal(t, "test_address", cfg.Address)
		assert.Equal(t, 9000, cfg.StoreInterval)
		assert.Equal(t, "my_file.data", cfg.FileStoragePath)
		assert.False(t, cfg.Restore)
		fmt.Println(cfg)
	}
	cmd := NewServerCmd()
	cmd.SetArgs([]string{"", "-a", "test_address", "-i", "9000", "-f", "my_file.data", "-r=False"})
	err := cmd.Execute()
	assert.NoError(t, err)
}

func Test_CanGetValuesFromEnv(t *testing.T) {
	_ = os.Setenv("ADDRESS", "test_address")
	_ = os.Setenv("STORE_INTERVAL", "1234")
	_ = os.Setenv("FILE_STORAGE_PATH", "asdf.data")
	_ = os.Setenv("RESTORE", "false")
	runServer = func(cfg *config.Config) {
		assert.Equal(t, "test_address", cfg.Address)
		assert.Equal(t, 1234, cfg.StoreInterval)
		assert.Equal(t, "asdf.data", cfg.FileStoragePath)
		assert.Equal(t, false, cfg.Restore)
	}
	cmd := NewServerCmd()
	err := cmd.Execute()
	assert.NoError(t, err)
}
