package password

import "gurms/internal/infra/property"

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

	return &PasswordManager{
		adminPasswordEncodingAlgorithm: adminPasswordEncodingAlgorithm,
		userPasswordEncodingAlgorithm:  userPasswordEncodingAlgorithm,
	}
}
