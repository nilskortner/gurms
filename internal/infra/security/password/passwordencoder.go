package password

type PasswordEncoder interface {
	hashPassword(password string) ([]byte, error)
	checkPassword(
		password string, hashedPassword []byte) bool
}
