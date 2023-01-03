package helper

import "golang.org/x/crypto/bcrypt"

func GeneratePassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	return string(hash), err
}
