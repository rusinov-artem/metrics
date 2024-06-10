package main

import (
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
