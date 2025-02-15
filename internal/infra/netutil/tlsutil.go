package netutil

import (
	"crypto/tls"
	"fmt"
	"gurms/internal/infra/property/env/common"
)

func CreateTlsConfig(tlsProperties *common.TlsProperties, forServer bool) *tls.Config {
	if !tlsProperties.Enabled {
		return nil
	}
	var tlsConfig *tls.Config
	if forServer {
		tlsConfig = getClientTlsConfig()
	} else {
		tlsConfig = createServerTlsConfig(tlsProperties)
	}
	return tlsConfig
}

func getClientTlsConfig() *tls.Config {
	trustManager
	if tls
}

func createServerTlsConfig(tlsProperties *common.TlsProperties) *tls.Config {
	cert, err := TlsCertificate()
	config := &tls.Config{}
	trustManager(getTrustManagerFactory())
	if tlsProperties.EnabledProtocols != nil {
		config.prot
	}
	if tlsProperties.ClientAuth != nil {
		config.ClientAuth = tlsProperties.ClientAuth
	}


	return config
}

func TlsCertificate(certPath, keyPath string) (tls.Certificate, error) {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("failed to load keystore: %w", err)
	}
	return cert, nil
}

// func TlsCA(caPath string) (*x509.CertPool, error) {
// 	certPool := x509.NewCertPool()
// 	caCert, err := ioutil.ReadFile(caPath)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to load truststore: %w", err)
// 	}
// 	if !certPool.AppendCertsFromPEM(caCert) {
// 		return nil, fmt.Errorf("failed to append CA certificates to truststore")
// 	}
// 	return certPool, nil
// }
