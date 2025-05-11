package password

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

const (
	CHARSET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+"
)

func GenerateSecret(length int) (string, error) {
	password := make([]byte, length)
	charsetLength := big.NewInt(int64(len(CHARSET)))
	for i := range length {
		index, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", fmt.Errorf("error whiel generating random indexes: %v", err)
		}
		password[i] = CHARSET[index.Int64()]
	}
	return string(password), nil

}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
