package util

import (
	"regexp"
)

// regex pattern to match IPv4 addresses
var ipv4Regex = regexp.MustCompile(`^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`)

// IsValidIPv4 checks if the given string is a valid IPv4 address.
func IsValidIPv4(ip string) bool {
	return ipv4Regex.MatchString(ip)
}
