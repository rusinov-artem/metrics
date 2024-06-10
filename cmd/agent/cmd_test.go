package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_CanHandleCommandLineArgs(t *testing.T) {
	runAgent = func(cfg config) {
		assert.Equal(t, "test_address", cfg.address)
		assert.Equal(t, 60*time.Second, cfg.pollInterval)
		assert.Equal(t, 90*time.Second, cfg.reportInterval)
	}
	cmd := NewAgent()
	cmd.SetArgs([]string{"", "-a", "test_address", "-p", "60s", "-r", "90s"})
	err := cmd.Execute()
	assert.NoError(t, err)
}
