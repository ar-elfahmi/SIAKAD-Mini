package service

import (
	"time"

	"github.com/ar-elfahmi/SIAKAD-Mini/app/model"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	JWTSecret string
}

func (s *AuthService) GenerateToken(user *model.User) (string, int, error) {
	expiresIn := 3600

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.JWTSecret))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expiresIn, nil
}
