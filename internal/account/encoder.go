package account

import (
	"golang.org/x/crypto/bcrypt"
)

// EncodePassword encode a given password with bcrypt and default cost
func EncodePassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// MatchesPassword returns a nil if matched hashedPassword and raw password, otherwise returns a error
func MatchesPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
