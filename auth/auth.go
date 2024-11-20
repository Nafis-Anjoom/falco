package auth

import (
	"fmt"
	"strconv"
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

func (as *AuthService) PasswordMatches(password string, hashedPassword []byte) bool {
    err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
    if err != nil {
        return false
    }

    return true
}

func (as *AuthService) NewToken(userId int64) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "falco",
		Subject:   strconv.FormatInt(userId, 10),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		NotBefore: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    fmt.Println(as.jwtSecret)
	return token.SignedString(as.jwtSecret)
}

func (as *AuthService) VerifyToken(tokenString string) (int, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return false, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}

		return as.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return -1, err
	}

    userId, _ := strconv.ParseInt(claims.Subject, 10, 64)
	return int(userId), nil
}
