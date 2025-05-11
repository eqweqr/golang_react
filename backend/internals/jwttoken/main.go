package jwttoken

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GetToken(bearer string) string {
	return strings.Split(bearer, " ")[1]
}

func ParseToken(tokenString string, secretKey string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("%v when parsing token", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

func CreateToken(secretKey string, userid int, username string, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userid":   userid,
			"exp":      time.Now().Add(time.Hour * 168).Unix(),
			"role":     role,
			"username": username,
		})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GetRoles(token *jwt.Token) string {

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if r, ok := claims["role"].(string); ok {
			return r
		}
	}
	return ""
}
