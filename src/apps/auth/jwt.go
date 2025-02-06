package auth

import (
	"errors"
	"socious-id/src/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

type Claims struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Refresh bool   `json:"refresh"`
	jwt.RegisteredClaims
}

func NewToken(id, email string) (*Token, error) {
	accessToken, err := GenerateToken(id, email, false)
	if err != nil {
		return nil, err
	}
	refreshToken, err := GenerateToken(id, email, true)
	if err != nil {
		return nil, err
	}
	return &Token{accessToken, refreshToken, "Bearer"}, nil

}

func GenerateToken(id, email string, refresh bool) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		ID:      id,
		Email:   email,
		Refresh: refresh,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Config.Secret))
}

func VerifyToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config.Secret), nil
	})
	if err != nil {
		return nil, err
	} else if !token.Valid {
		return nil, errors.New("invalid token")
	} else if claims, ok := token.Claims.(*Claims); ok {
		return claims, nil
	}

	return nil, errors.New("unknown claims type, cannot proceed")
}
