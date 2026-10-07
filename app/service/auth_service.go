package service

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/ar-elfahmi/SIAKAD-Mini/app/model"
	"github.com/ar-elfahmi/SIAKAD-Mini/app/repository"
)

type AuthService struct {
	Repository *repository.AuthRepository
	JWTSecret  string
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
func (s *AuthService) GetUserByID(
	ctx context.Context,
	id int,
) (*model.User, error) {
	return s.Repository.GetUserByID(ctx, id)
}
