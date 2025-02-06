package auth

import (
	"golang.org/x/crypto/bcrypt"
)

type RegisterForm struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Username  *string `json:"username"`
	Email     string  `json:"email" validate:"required,email"`
	Password  *string `json:"password"`
}

type LoginForm struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,min=8"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func GenerateFullTokens(id, email string) (map[string]any, error) {
	accessToken, err := GenerateToken(id, email, false)
	if err != nil {
		return nil, err
	}
	refreshToken, err := GenerateToken(id, email, true)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_type":    "Bearer",
	}, nil
}
