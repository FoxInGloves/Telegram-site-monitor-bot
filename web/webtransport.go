package web

import (
	"context"
	"errors"
	"golang.org/x/net/proxy"
	"net"
	"net/http"
	"net/url"
	"strings"
)

func GetWebTransport(proxyAddress string) (*http.Transport, error) {
	isValidProxy := isValidUrlOrAddr(proxyAddress)
	if !isValidProxy {
		return nil, errors.New("invalid proxy address")
	}

	if !strings.Contains(proxyAddress, "://") {
		proxyAddress = "http://" + proxyAddress
	}

	parsedURL, err := url.Parse(proxyAddress)
	if err != nil {
		return nil, err
	}

	var transport *http.Transport

	switch parsedURL.Scheme {
	case "http", "https":
		transport = &http.Transport{
			Proxy: http.ProxyURL(parsedURL),
		}
	case "socks5", "socks5h":
		dialer, proxyErr := proxy.FromURL(parsedURL, proxy.Direct)
		if proxyErr != nil {
			return nil, err
		}
		transport = &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
			},
		}
	default:
		return nil, errors.New("unsupported proxy scheme: " + parsedURL.Scheme)
	}

	return transport, nil
}
