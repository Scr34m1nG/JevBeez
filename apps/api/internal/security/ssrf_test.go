package security

import (
	"net"
	"testing"
)

func TestIsPrivateIP(t *testing.T) {
	tests := []struct {
		ip       string
		expected bool
	}{
		{"0.0.0.0", true},
		{"127.0.0.1", true},
		{"10.0.0.1", true},
		{"172.16.0.1", true},
		{"172.20.0.5", true},
		{"192.168.1.1", true},
		{"169.254.169.254", true},
		{"::1", true},
		{"::ffff:127.0.0.1", true},
		{"::ffff:10.0.0.1", true},
		{"::ffff:169.254.169.254", true},
		{"::ffff:8.8.8.8", false},
		{"8.8.8.8", false},
		{"1.1.1.1", false},
	}

	for _, tc := range tests {
		t.Run(tc.ip, func(t *testing.T) {
			ip := net.ParseIP(tc.ip)
			got := IsPrivateIP(ip)
			if got != tc.expected {
				t.Errorf("IsPrivateIP(%s) = %v, expected %v", tc.ip, got, tc.expected)
			}
		})
	}
}
