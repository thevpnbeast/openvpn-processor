package options

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetOpenvpnProcessorOptions(t *testing.T) {
	opts := GetOpenvpnProcessorOptions()
	assert.NotNil(t, opts)
}
