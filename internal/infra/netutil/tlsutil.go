package netutil

import (
	"crypto/tls"
	"gurms/internal/infra/property/env/common"
)

func CreateTlsConfig(tlsProperties *common.TlsProperties, forServer bool) *tls.Config {
	if !tlsProperties.Enabled {
		return nil
	}
	var tlsConfig *tls.Config
	if forServer {
		tlsConfig = createServerTlsConfig(tlsProperties)
	} else {
		tlsConfig = createClientTlsConfig(tlsProperties)
	}
	return tlsConfig
}

// TODO: make longlived 10year CA and shortterm 1year CAs that renew on expiry
func createServerTlsConfig(tlsProperties *common.TlsProperties) *tls.Config {
	config := &tls.Config{
		Certificates: []tls.Certificate{*tlsProperties.Certificate},
		ClientCAs:    tlsProperties.CA,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		MinVersion:   tls.VersionTLS13,
		MaxVersion:   tls.VersionTLS13,
	}
	return config
}

func createClientTlsConfig(tlsProperties *common.TlsProperties) *tls.Config {
	config := &tls.Config{
		RootCAs:    tlsProperties.CA,
		MinVersion: tls.VersionTLS13,
		MaxVersion: tls.VersionTLS13,
	}
	return config
}

// func tlsCertificate(certPath, keyPath string) (tls.Certificate, error) {
// 	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
// 	if err != nil {
// 		return tls.Certificate{}, fmt.Errorf("failed to load keystore: %w", err)
// 	}
// 	return cert, nil
// }

// func tlsCA(caPath string) (*x509.CertPool, error) {
// 	clientCAs := x509.NewCertPool()
// 	caCert, err := os.ReadFile(caPath)
// 	if err != nil {
// 		return clientCAs, err
// 	}
// 	clientCAs.AppendCertsFromPEM(caCert)

// 	return clientCAs, nil
// }

// cert, err := tlsCertificate("server.crt", "server.key")
// if err != nil {
// 	return &tls.Config{}, err
// }
// clientCAs, err := tlsCA("ca.pem")
// if err != nil {
// 	return &tls.Config{}, err
// }
