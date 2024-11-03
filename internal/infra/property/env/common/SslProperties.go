package common

type SslProperties struct {
	enabled          bool
	clientAuth       ClienthAuth
	ciphers          []string
	enabledProtocols []string

	keyAlias         string
	keyPassword      string
	keyStore         string
	keyStorePassword string
	keyStoreType     string
	keyStoreProvider string

	trustStore         string
	trustStorePassword string
	trustStoreType     string
	trustStoreProvider string
	protocol           string
}
