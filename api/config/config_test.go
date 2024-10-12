package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	tests := []struct {
		name        string
		mockEnvFunc func() func()
		fail        bool
	}{
		{
			"error due to bad port",
			func() func() {
				os.Setenv("VITTA_PORT", "invalid-2839")

				return func() {
					os.Unsetenv("VITTA_PORT")
				}
			},
			true,
		},
		{
			"success",
			func() func() {
				os.Setenv("VITTA_PORT", "15002")

				return func() {
					os.Unsetenv("VITTA_PORT")
				}
			},
			false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			defer tc.mockEnvFunc()()

			cfg, err := New()

			if tc.fail {
				require.Error(t, err)
				assert.Nil(t, cfg)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, cfg)
			}
		})
	}
}
