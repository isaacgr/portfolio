package client

import (
	"crypto/tls"
	"math"
	"net/http"
	"time"
)

type HttpClientOptions struct {
	RequestTimeout int
}

func HttpClient(
	httpClientOptions *HttpClientOptions,
) *http.Client {
	transport := &http.Transport{
		DisableKeepAlives:   false,
		MaxIdleConns:        0,
		MaxIdleConnsPerHost: math.MaxInt64,
		IdleConnTimeout:     300 * time.Second,
		TLSClientConfig: &tls.Config{
			// TODO: Add certificate verification
			InsecureSkipVerify: true,
		},
	}
	return &http.Client{
		Transport: transport,
		Timeout: time.Duration(
			httpClientOptions.RequestTimeout,
		) * time.Second,
	}
}
