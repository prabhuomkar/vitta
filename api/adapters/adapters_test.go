package adapters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name       string
		expAdapter Adapter
	}{
		{"icici", newICICI("category")},
		{"empty", nil},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			adapter := New(tc.name, "category", nil)
			assert.Equal(t, tc.expAdapter, adapter)
		})
	}
}
