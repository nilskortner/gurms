package security

type SecurityProperties struct {
	password  *PasswordProperties
	blocklist *BlocklistProperties
}

func NewSecurityProperties() *SecurityProperties {
	return &SecurityProperties{
		password:  NewPasswordProperties(),
		blocklist: NewBlocklistProperties(),
	}
}
