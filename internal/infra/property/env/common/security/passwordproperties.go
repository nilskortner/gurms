package security

import "gurms/internal/infra/property/constants"

type PasswordProperties struct {
	initialRootPassword            string
	userPasswordEncodingAlgorithm  string
	adminPasswordEncodingAlgorithm string
}

func NewPasswordProperties() *PasswordProperties {
	return &PasswordProperties{
		initialRootPassword:            "gurms",
		userPasswordEncodingAlgorithm:  constants.SALTED_SHA256,
		adminPasswordEncodingAlgorithm: constants.BCRYPT,
	}
}
