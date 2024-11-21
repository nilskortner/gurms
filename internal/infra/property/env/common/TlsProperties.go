package common

import "crypto/tls"

type TlsProperties struct {
	Enabled          bool
	ClientAuth       tls.ClientAuthType
	Ciphers          []string
	EnabledProtocols []string

	KeyAlias         string
	KeyPassword      string
	KeyStore         string
	KeyStorePassword string
	KeyStoreType     string
	KeyStoreProvider string

	TrustStore         string
	TrustStorePassword string
	TrustStoreType     string
	TrustStoreProvider string
	Protocol           string
}
