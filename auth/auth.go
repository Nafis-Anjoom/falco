package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	jwtSecret          []byte
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

func NewAuthService(jwtSecret string) *AuthService {
	return &AuthService{
		jwtSecret:          []byte(jwtSecret),
		accessTokenExpiry:  time.Hour * 24,
		refreshTokenExpiry: time.Hour * 24 * 7,
	}
}

func (as *AuthService) HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 12)
}

func (as *AuthService) NewToken(userId int) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "falco",
		Subject:   string(userId),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(as.refreshTokenExpiry)),
		NotBefore: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	return token.SignedString(as.jwtSecret)
}

func (as *AuthService) VerifyToken(tokenString string) (*jwt.RegisteredClaims, error) {
    claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return false, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}

        return as.jwtSecret, nil
	})

    if err != nil || !token.Valid {
        return nil, err
    }

    return claims, nil
}
