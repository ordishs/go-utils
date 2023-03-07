package utils

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIP(t *testing.T) {
	host, port, err := net.SplitHostPort(":8001")

	assert.Len(t, host, 0)
	assert.Equal(t, "8001", port)
	assert.Nil(t, err)
}
