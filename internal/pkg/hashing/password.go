package hashing

import (
	"example/auth/internal/pkg/models"

	"golang.org/x/crypto/bcrypt"
)

// typecase models.User to User
type User models.User

func Hash(password string) (string, error) {
	salt, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		return "", err
	}

	hashedPassword := string(salt) + password

	return hashedPassword, nil
}

func ComparePassword(userpassword string, dbpassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(dbpassword), []byte(userpassword))

	if err != nil {
		return false
	}

	return true
}
