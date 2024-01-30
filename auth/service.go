package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service interface {
	GenerateToken(userID int) (string, error)
}

type service struct {
	jwtSecret string
}

func NewService(jwtSecret string) *service {
	return &service{jwtSecret}
}

func (s *service) GenerateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id":    userID,
		"expired_at": time.Now().Add(time.Hour * 24 * 100).Format(time.RFC822),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return signedToken, err
	}

	return signedToken, nil
}
