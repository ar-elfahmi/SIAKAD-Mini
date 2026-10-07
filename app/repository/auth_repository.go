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
		`SELECT
    u.id,
    u.email,
    u.role
		FROM users u
		LEFT JOIN students s ON s.user_id = u.id
		WHERE u.email = $1
  AND (
      u.role = 'admin'
      OR s.deleted_at IS NULL
  )`,
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
func (r *AuthRepository) GetUserByID(ctx context.Context, id int) (*model.User, error) {
	var user model.User
	var student model.StudentProfile

	err := r.Pool.QueryRow(ctx, `
		SELECT
			u.id,
			u.email,
			u.role,
			s.nim,
			s.nama,
			s.prodi,
			s.angkatan
		FROM users u
		LEFT JOIN students s
			ON s.user_id = u.id
			AND s.deleted_at IS NULL
		WHERE u.id = $1
	`, id).Scan(
		&user.ID,
		&user.Email,
		&user.Role,
		&student.NIM,
		&student.Nama,
		&student.Prodi,
		&student.Angkatan,
	)

	if err != nil {
		return nil, err
	}

	if user.Role == "mahasiswa" {
		user.Students = &student
	}

	return &user, nil
}
