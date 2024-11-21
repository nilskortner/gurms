package netutil

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"gurms/internal/infra/property/env/common"
	"io/ioutil"
)

func CreateTlsConfig(tlsProperties *common.TlsProperties, forServer bool) *tls.Config {
	if !tlsProperties.Enabled {
		return nil
	}
	if forServer {

	} else {

	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    tlsProperties.TrustStore,
	}

	if requireClientCert {
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
	}
	return tlsConfig
}

func TlsCertificate(certPath, keyPath string) (tls.Certificate, error) {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("failed to load keystore: %w", err)
	}
	return cert, nil
}

func TlsCA(caPath string) (*x509.CertPool, error) {
	certPool := x509.NewCertPool()
	caCert, err := ioutil.ReadFile(caPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load truststore: %w", err)
	}
	if !certPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to append CA certificates to truststore")
	}
	return certPool, nil
}
