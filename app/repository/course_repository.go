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
	search string,
	avilable bool,) ([]model.Course, error) {
	courses := []model.Course{}

	rows, err := r.Pool.Query( 
		ctx,
		`SELECT c.id, c.kode_mk, c.nama_mk, c.sks, c.semester, c.kuota,
			COUNT(e.course_id) AS terisi,
			c.kuota - COUNT(e.course_id) AS sisa_kuota
		FROM courses c
		LEFT JOIN enrollments e ON e.course_id = c.id
		WHERE c.semester = $1
		AND (
			c.kode_mk ILIKE $2
			OR c.nama_mk ILIKE $2
			)
		GROUP BY c.id
		HAVING (
		NOT $3
		OR COUNT(e.id) < c.kuota)`,
		semester,
		"%"+search+"%",
		avilable,
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
			&course.Terisi,
			&course.SisaKuota,
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
