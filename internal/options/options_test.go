package options

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetVpnbeastServiceOptions(t *testing.T) {
	opts := GetVpnbeastServiceOptions()
	assert.NotNil(t, opts)
}
