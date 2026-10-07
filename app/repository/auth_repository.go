package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ar-elfahmi/SIAKAD-Mini/app/model"
)

type AuthRepository struct {
	Pool *pgxpool.Pool
}

func (r *AuthRepository) FindUserByEmail(
	ctx context.Context,
	email string,
) (*model.User, string, error) {

	var user model.User
	var passwordHash string

	err := r.Pool.QueryRow(
		ctx,
		`SELECT id, email, password, role
		FROM users
		WHERE email = $1`,
		email,
	).Scan(
		&user.ID,
		&user.Email,
		&passwordHash,
		&user.Role,
	)

	if err != nil {
		return nil, "", err
	}

	return &user, passwordHash, nil
}
