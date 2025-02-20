package common

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

type TlsProperties struct {
	Enabled    bool
	ClientAuth tls.ClientAuthType

	Certificate *tls.Certificate
	CA          *x509.CertPool
}

// TODO: create certificates for property init
func NewTlsProperties() *TlsProperties {
	cert, err := tlsCertificate("", "") //<-
	if err != nil {
		panic("could not create TlsProperties missing Path")
	}
	ca, err := tlsCA("") //<-
	if err != nil {
		panic("could not create TlsProperties missing Path")
	}
	return &TlsProperties{
		Enabled:     true,
		ClientAuth:  tls.RequireAndVerifyClientCert,
		Certificate: &cert,
		CA:          ca,
	}
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
