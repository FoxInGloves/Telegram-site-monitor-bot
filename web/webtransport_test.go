package web

import (
	"strings"
	"testing"
)

func TestGetWebTransport_Http_HttpTransport(t *testing.T) {
	transport, err := GetWebTransport("http://127.0.0.1")
	if err != nil {
		t.Fatalf("error getting web transport: %v", err)
	}
	if transport == nil {
		t.Error("transport is nil")
	}
	if transport.Proxy == nil {
		t.Error("proxy is nil")
	}
	if transport.DialContext != nil {
		t.Error("dial context is not nil")
	}
}

func TestGetWebTransport_Nothing_HttpTransport(t *testing.T) {
	transport, err := GetWebTransport("127.0.0.1:90")
	if err != nil {
		t.Fatalf("error getting web transport: %v", err)
	}
	if transport == nil {
		t.Error("transport is nil")
	}
	if transport.Proxy == nil {
		t.Error("proxy is nil")
	}
	if transport.DialContext != nil {
		t.Error("dial context is not nil")
	}
}

func TestGetWebTransport_Socks5_Socks5Transport(t *testing.T) {
	transport, err := GetWebTransport("socks5://127.0.0.1")
	if err != nil {
		t.Fatalf("error getting web transport: %v", err)
	}
	if transport == nil {
		t.Error("transport is nil")
	}
	if transport.Proxy != nil {
		t.Error("proxy is nil")
	}
	if transport.DialContext == nil {
		t.Error("dial context is nil")
	}
}

func TestGetWebTransport_Errors(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectError   bool
		errorContains string
	}{
		{
			name:        "Invalid URL syntax",
			input:       "http://[::1",
			expectError: true,
		},
		{
			name:          "Unsupported scheme",
			input:         "ftp://example.com:3128",
			expectError:   true,
			errorContains: "unsupported proxy scheme",
		},
		{
			name:        "Broken SOCKS5 URL",
			input:       "socks5://invalid::proxy",
			expectError: true,
		},
		{
			name:        "Empty scheme and invalid format",
			input:       "::::::",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transport, err := GetWebTransport(tt.input)

			if tt.expectError && err == nil {
				t.Errorf("expected error but got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("did not expect error but got: %v", err)
			}
			if err != nil && tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
				t.Errorf("expected error to contain %q, got %q", tt.errorContains, err.Error())
			}
			if !tt.expectError && transport == nil {
				t.Errorf("expected non-nil transport, got nil")
			}
		})
	}
}
