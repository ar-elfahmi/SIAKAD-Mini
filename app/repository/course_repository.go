package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ar-elfahmi/SIAKAD-Mini/app/model"
)

type CourseRepository struct {
	Pool *pgxpool.Pool
}

func (r *CourseRepository) GetAllCourses(
	ctx context.Context, 
	semester int,
	search string,) ([]model.Course, error) {
	courses := []model.Course{}

	rows, err := r.Pool.Query( 
		ctx,
		`SELECT id, kode_mk, nama_mk, sks, semester, kuota 
		FROM courses
		WHERE semester = $1
		AND (
			kode_mk ILIKE $2
			OR nama_mk ILIKE $2
			)`,
		semester,
		"%"+search+"%",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var course model.Course

		err := rows.Scan(
			&course.ID,
			&course.KodeMK,
			&course.NamaMK,
			&course.SKS,
			&course.Semester,
			&course.Kuota,
		)

		if err != nil {
			return nil, err
		}

		courses = append(courses, course)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return courses, nil
}
