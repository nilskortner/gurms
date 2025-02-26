package security

type SecurityProperties struct {
	Password  *PasswordProperties
	Blocklist *BlocklistProperties
}

func NewSecurityProperties() *SecurityProperties {
	return &SecurityProperties{
		Password:  NewPasswordProperties(),
		Blocklist: NewBlocklistProperties(),
	}
}
