package web

import "testing"

func TestIsValidUrlOrAddr(t *testing.T) {
	tests := []struct {
		urlOrAddr string
		expected  bool
	}{
		{"http://example.com", true},
		{"https://google.com", true},
		{"socks5://127.0.0.1", true},
		{"socks5h://122.2.2.1", true},
		{"ftp://server.local", true},
		{"example", false},
		{"1234567890", false},
		{"localhost:8080", false},
		{"127.0.0.1", false},
		{"sub.domain.com:8080", true},
		{"127.0.0.1:8080", true},
		{"example.com", false},
		{"google.com", false},
		{"127.0.0.1:8080:8081", false},
		{"127.0.0.1:80800", false},
		{"134.45.234.1:-1", false},
	}

	for _, tt := range tests {
		if isValidUrlOrAddr(tt.urlOrAddr) != tt.expected {
			t.Errorf("expected %v, received %v (%v)", tt.expected, isValidUrlOrAddr(tt.urlOrAddr), tt.urlOrAddr)
		}
	}
}
