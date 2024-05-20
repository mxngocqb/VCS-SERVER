package util

import "testing"

func TestIsValidIPv4(t *testing.T) {
	testCases := []struct {
		ip       string
		expected bool
	}{
		{"192.168.0.1", true},    // valid IPv4 address
		{"256.256.256.256", false}, // invalid IPv4 address (out of range)
		{"0.0.0.0", true},         // valid IPv4 address
		{"127.0.0.1", true},       // valid IPv4 address (localhost)
		{"1.2.3.4", true},         // valid IPv4 address
		{"123.456.789.0", false},  // invalid IPv4 address (out of range)
		{"abc.def.ghi.jkl", false}, // invalid format
		{"", false},               // empty string
	}

	for _, tc := range testCases {
		t.Run(tc.ip, func(t *testing.T) {
			result := IsValidIPv4(tc.ip)
			if result != tc.expected {
				t.Errorf("Expected %v for %s, but got %v", tc.expected, tc.ip, result)
			}
		})
	}
}
