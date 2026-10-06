package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ar-elfahmi/SIAKAD-Mini/app/model"
)

type StudentRepository struct {
	Pool *pgxpool.Pool
}

func (r *StudentRepository) GetAllStudents(
	ctx context.Context,
	page, perPage int,
	prodi string,
	angkatan int,
	search, sort string,
) ([]model.Student, int, error) {
	offset := (page - 1) * perPage

	orderBy := "nama ASC"
	switch sort {
	case "nama":
		orderBy = "nama ASC"
	case "-ipk_terakhir":
		orderBy = "ipk_terakhir DESC"
	}

	rows, err := r.Pool.Query(
		ctx,
		`SELECT id, nim, nama, prodi, angkatan, ipk_terakhir,
			COUNT(*) OVER() AS total
		FROM students
		WHERE deleted_at IS NULL
			AND ($1 = '' OR prodi = $1)
			AND ($2 = 0 OR angkatan = $2)
			AND ($3 = '' OR nim ILIKE '%' || $3 || '%' OR nama ILIKE '%' || $3 || '%')
		ORDER BY `+orderBy+`
		LIMIT $4 OFFSET $5`,
		prodi,
		angkatan,
		search,
		perPage,
		offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	students := []model.Student{}
	total := 0

	for rows.Next() {
		var s model.Student
		var rowTotal int
		err := rows.Scan(
			&s.ID,
			&s.NIM,
			&s.Nama,
			&s.Prodi,
			&s.Angkatan,
			&s.IPKTerakhir,
			&rowTotal,
		)
		if err != nil {
			return nil, 0, err
		}
		students = append(students, s)
		total = rowTotal
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return students, total, nil
}
