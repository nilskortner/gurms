package netutil

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"gurms/internal/infra/property/env/common"
	"os"
)

func CreateTlsConfig(tlsProperties *common.TlsProperties, forServer bool) (*tls.Config, error) {
	if !tlsProperties.Enabled {
		return nil, nil
	}
	var tlsConfig *tls.Config
	var err error
	if forServer {
		tlsConfig, err = createServerTlsConfig(tlsProperties)
	} else {
		tlsConfig, err = createClientTlsConfig(tlsProperties)
	}
	return tlsConfig, err
}

// TODO: make longlived 10year CA and shortterm 1year CAs that renew on expiry
func createServerTlsConfig(tlsProperties *common.TlsProperties) (*tls.Config, error) {
	cert, err := tlsCertificate("server.crt", "server.key")
	if err != nil {
		return &tls.Config{}, err
	}
	clientCAs, err := tlsCA("ca.pem")
	if err != nil {
		return &tls.Config{}, err
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    clientCAs,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		MinVersion:   tls.VersionTLS13,
		MaxVersion:   tls.VersionTLS13,
	}
	return config, nil
}

func createClientTlsConfig(tlsProperties *common.TlsProperties) (*tls.Config, error) {
	rootCAs, err := tlsCA("ca.pem")
	if err != nil {
		return &tls.Config{}, err
	}
	config := &tls.Config{
		RootCAs:    rootCAs,
		MinVersion: tls.VersionTLS13,
		MaxVersion: tls.VersionTLS13,
	}
	return config, nil
}

func tlsCertificate(certPath, keyPath string) (tls.Certificate, error) {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("failed to load keystore: %w", err)
	}
	return cert, nil
}

func tlsCA(caPath string) (*x509.CertPool, error) {
	clientCAs := x509.NewCertPool()
	caCert, err := os.ReadFile(caPath)
	if err != nil {
		return clientCAs, err
	}
	clientCAs.AppendCertsFromPEM(caCert)

	return clientCAs, nil
}
