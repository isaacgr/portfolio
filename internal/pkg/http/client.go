package http

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/isaacgr/portfolio/internal/pkg/logging"
)

var log = logging.GetLogger("internal.http", false)

// For security puporses only load CA Certs from this path
const caCertPath = "/var/lib/portfolio/ca-certs"

type HttpClientOptions struct {
	RequestTimeout int
	TLSConfig      *tls.Config
}

func loadCaCertsFromPath(dir string) (*x509.CertPool, error) {
	pool := x509.NewCertPool()

	err := filepath.Walk(
		dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}

			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			var block *pem.Block
			rest := data
			for {
				block, rest = pem.Decode(rest)
				if block == nil {
					break
				}
				if block.Type == "CERTIFICATE" {
					cert, err := x509.ParseCertificate(block.Bytes)
					if err != nil {
						return err
					}
					pool.AddCert(cert)
				}
			}

			return nil
		})

	if err != nil {
		return nil, err
	}
	return pool, nil
}

func httpClient(
	httpClientOptions *HttpClientOptions,
) *http.Client {
	transport := &http.Transport{
		DisableKeepAlives:   false,
		MaxIdleConns:        0,
		MaxIdleConnsPerHost: math.MaxInt64,
		IdleConnTimeout:     300 * time.Second,
		TLSClientConfig:     httpClientOptions.TLSConfig,
	}
	return &http.Client{
		Transport: transport,
		Timeout: time.Duration(
			httpClientOptions.RequestTimeout,
		) * time.Second,
	}
}

func NewHttpClient() *http.Client {

	caCertPool, err := loadCaCertsFromPath(caCertPath)

	if err != nil {
		log.Error("Unable to load CA Certificates. ", "Error", err)
	}

	tlsConfig := &tls.Config{
		RootCAs: caCertPool,
	}

	client := httpClient(&HttpClientOptions{
		RequestTimeout: 30,
		TLSConfig:      tlsConfig,
	})

	return client
}
