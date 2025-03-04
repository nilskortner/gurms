package password

import (
	"gurms/internal/infra/property"
	"gurms/internal/infra/property/constants"
)

type PasswordManager struct {
	adminPasswordEncodingAlgorithm string
	userPasswordEncodingAlgorithm  string
	bCryptPasswordEncoder          *BCryptPasswordEncoder
	sha256PasswordEncoder          *SaltedSha256PasswordEncoder
}

func NewPasswordManager(propertiesManager *property.GurmsPropertiesManager) *PasswordManager {
	gurmsProperties := propertiesManager.LocalGurmsProperties
	adminPasswordEncodingAlgorithm := gurmsProperties.Security.Password.GetAdmingPasswordEncodingAlgorithm()
	userPasswordEncodingAlgorithm := gurmsProperties.Security.Password.GetUserPasswordEncodingAlgorithm()

	var bCryptPasswordEncoder *BCryptPasswordEncoder
	if adminPasswordEncodingAlgorithm == constants.BCRYPT ||
		userPasswordEncodingAlgorithm == constants.BCRYPT {
		bCryptPasswordEncoder = &BCryptPasswordEncoder{}
	}
	var sha256PasswordEncoder *SaltedSha256PasswordEncoder
	if adminPasswordEncodingAlgorithm == constants.SALTED_SHA256 ||
		userPasswordEncodingAlgorithm == constants.SALTED_SHA256 {
		sha256PasswordEncoder = &SaltedSha256PasswordEncoder{}
	}

	return &PasswordManager{
		adminPasswordEncodingAlgorithm: adminPasswordEncodingAlgorithm,
		userPasswordEncodingAlgorithm:  userPasswordEncodingAlgorithm,
		bCryptPasswordEncoder:          bCryptPasswordEncoder,
		sha256PasswordEncoder:          sha256PasswordEncoder,
	}
}
