package password

import "golang.org/x/crypto/bcrypt"

type BCryptPasswordEncoder struct {
}

func (b *BCryptPasswordEncoder) hashPassword(password string) ([]byte, error) {
	// Generate a bcrypt hash with cost factor 10
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func (b *BCryptPasswordEncoder) checkPassword(password string, hashedPassword []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	return err == nil // Returns true if password matches
}
