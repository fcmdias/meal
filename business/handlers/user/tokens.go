package user

import (
	"github.com/fcmdias/meal/business/web/auth"
)

type TokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func generateToken(email, username string) (string, error) {
	tokenString, err := auth.GenerateJWT(email, username)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
