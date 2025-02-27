package hc

import (
	"net"
	"net/http"
	"time"
)

var (
	httpClient *http.Client
)

func init() {
	dialer := &net.Dialer{
		Timeout:   2 * time.Minute,
		KeepAlive: 2 * time.Minute,
	}
	httpClient = &http.Client{
		Transport: &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			DialContext:           dialer.DialContext,
			MaxIdleConns:          1000,
			MaxIdleConnsPerHost:   300,
			IdleConnTimeout:       2 * time.Minute,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
}

func Client() *http.Client {
	return httpClient
}

func Do(req *http.Request) (*http.Response, error) {
	return httpClient.Do(req)
}
