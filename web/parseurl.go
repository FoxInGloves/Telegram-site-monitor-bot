package web

import (
	"net"
	"net/url"
	"strconv"
	"strings"
)

func isValidUrlOrAddr(input string) bool {
	if u, err := url.Parse(input); err == nil && u.Scheme != "" && u.Host != "" {
		return true
	}

	host, port, err := net.SplitHostPort(input)
	if err != nil {
		return false
	}

	if portNum, err := strconv.Atoi(port); err != nil || portNum < 0 || portNum > 65535 {
		return false
	}

	if ip := net.ParseIP(host); ip != nil {
		return true
	}

	if strings.Contains(host, ".") {
		return true
	}

	return false
}
