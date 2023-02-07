package utils

import (
	"net"
	"testing"
)

func TestIP(t *testing.T) {
	host, port, err := net.SplitHostPort(":8001")

	t.Logf("Host %s (%d)", host, len(host))
	t.Logf("Port %s", port)
	t.Logf("err %v", err)
}
