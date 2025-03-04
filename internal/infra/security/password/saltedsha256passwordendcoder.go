package password

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"gurms/internal/infra/collection"
)

const SALT_SIZE_BYTES = 8

type SaltedSha256PasswordEncoder struct {
}

// TODO: testing

func generateSalt() ([]byte, error) {
	salt := make([]byte, SALT_SIZE_BYTES)
	_, err := rand.Read(salt)
	return salt, err
}

func (s *SaltedSha256PasswordEncoder) hashPassword(password string) ([]byte, error) {
	salt := make([]byte, SALT_SIZE_BYTES)
	_, err := rand.Read(salt)
	if err != nil {
		return salt, err
	}

	hasher := sha256.New()
	hasher.Write(salt)
	hasher.Write([]byte(password))
	sum := hasher.Sum(nil)

	return collection.Concat(salt, sum), nil
}

func (s *SaltedSha256PasswordEncoder) checkPassword(
	password string, saltedPasswordWithSalt []byte) bool {

	if len(saltedPasswordWithSalt) < SALT_SIZE_BYTES {
		return false
	}

	rawPasswordWithSalt := make([]byte, SALT_SIZE_BYTES+len(password))
	copy(rawPasswordWithSalt, saltedPasswordWithSalt[:SALT_SIZE_BYTES])
	copy(rawPasswordWithSalt[SALT_SIZE_BYTES:], []byte(password))

	hasher := sha256.New()
	hasher.Write(rawPasswordWithSalt)
	saltedPasswort := hasher.Sum(nil)

	return bytes.Equal(saltedPasswort, saltedPasswordWithSalt[SALT_SIZE_BYTES:])
}
