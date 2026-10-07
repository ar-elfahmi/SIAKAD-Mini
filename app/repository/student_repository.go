package repository

import (
	"context"
	"encoding/json"

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

func (r *StudentRepository) GetStudentByID(
	ctx context.Context,
	id int,
) (*model.StudentDetail, error) {
	var d model.StudentDetail
	var coursesJSON []byte

	err := r.Pool.QueryRow(
		ctx,
		`SELECT
			s.id, s.nim, s.nama, s.prodi, s.angkatan, s.ipk_terakhir,
			COALESCE(SUM(c.sks) FILTER (WHERE c.id IS NOT NULL), 0) AS total_sks,
			CASE
				WHEN s.ipk_terakhir >= 3.00 THEN 24
				WHEN s.ipk_terakhir >= 2.50 THEN 21
				WHEN s.ipk_terakhir >= 2.00 THEN 18
				ELSE 15
			END AS batas_sks,
			COALESCE(
				json_agg(
					json_build_object(
						'id', c.id,
						'kode_mk', c.kode_mk,
						'nama_mk', c.nama_mk,
						'sks', c.sks,
						'semester', c.semester,
						'tahun_akademik', e.tahun_akademik
					) ORDER BY c.semester, c.kode_mk
				) FILTER (WHERE c.id IS NOT NULL),
				'[]'::json
			) AS courses
		FROM students s
		LEFT JOIN enrollments e ON e.student_id = s.id
		LEFT JOIN courses c ON c.id = e.course_id
		WHERE s.id = $1 AND s.deleted_at IS NULL
		GROUP BY s.id`,
		id,
	).Scan(
		&d.ID,
		&d.NIM,
		&d.Nama,
		&d.Prodi,
		&d.Angkatan,
		&d.IPKTerakhir,
		&d.TotalSKS,
		&d.BatasSKS,
		&coursesJSON,
	)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(coursesJSON, &d.Courses); err != nil {
		return nil, err
	}
	if d.Courses == nil {
		d.Courses = []model.EnrolledCourse{}
	}

	return &d, nil
}
