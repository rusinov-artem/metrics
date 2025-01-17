package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rusinov-artem/metrics/cmd/server/config"
	"github.com/rusinov-artem/metrics/test"
)

func Test_CanHandleCommandLineArgs(t *testing.T) {
	m, destructFN := test.NewEnvManager()
	defer destructFN()
	m.Set("DATABASE_DSN", "")
	m.Set("KEY", "")
	runServer = func(cfg *config.Config) {
		assert.Equal(t, "test_address", cfg.Address)
		assert.Equal(t, 9000, cfg.StoreInterval)
		assert.Equal(t, "my_file.data", cfg.FileStoragePath)
		assert.Equal(t, "metrics-db:5432", cfg.DatabaseDsn)
		assert.False(t, cfg.Restore)
		assert.Equal(t, "some_key", cfg.Key)
		assert.Equal(t, "private.pem", cfg.CryptoKey)
		assert.Equal(t, "config.json", cfg.Config)
		fmt.Println(cfg)
	}
	cmd := NewServerCmd()
	cmd.SetArgs(
		[]string{"",
			"-a", "test_address",
			"-i", "9000",
			"-f", "my_file.data",
			"-r=False",
			"-d", "metrics-db:5432",
			"-k", "some_key",
			"--crypto-key", "private.pem",
			"--config", "config.json",
		},
	)
	err := cmd.Execute()
	assert.NoError(t, err)
}

func Test_CanGetValuesFromEnv(t *testing.T) {
	m, destructFN := test.NewEnvManager()
	defer destructFN()
	m.Set("ADDRESS", "test_address")
	m.Set("STORE_INTERVAL", "1234")
	m.Set("FILE_STORAGE_PATH", "asdf.data")
	m.Set("RESTORE", "false")
	m.Set("DATABASE_DSN", "metrics-db:6432")
	m.Set("KEY", "some_key_from_env")
	m.Set("CRYPTO_KEY", "private.pem")
	m.Set("CONFIG", "config.json")
	runServer = func(cfg *config.Config) {
		assert.Equal(t, "test_address", cfg.Address)
		assert.Equal(t, 1234, cfg.StoreInterval)
		assert.Equal(t, "asdf.data", cfg.FileStoragePath)
		assert.Equal(t, "metrics-db:6432", cfg.DatabaseDsn)
		assert.Equal(t, "some_key_from_env", cfg.Key)
		assert.Equal(t, false, cfg.Restore)
		assert.Equal(t, "private.pem", cfg.CryptoKey)
		assert.Equal(t, "config.json", cfg.Config)
	}
	cmd := NewServerCmd()
	err := cmd.Execute()
	assert.NoError(t, err)
}
