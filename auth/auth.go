package auth

import "golang.org/x/crypto/bcrypt"


type AuthService struct {
    jwtSecret string
    accessTokenExpiry string
    refreshTokenExpiry string
}

func NewAuthService(jwtSecret, accessTokenExpiry, refreshTokenExpiry string) *AuthService {
    return &AuthService{
        jwtSecret: jwtSecret,
        accessTokenExpiry: accessTokenExpiry,
        refreshTokenExpiry: refreshTokenExpiry,
    }
}

func (as *AuthService) HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 12)
}
