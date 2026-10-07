package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ar-elfahmi/SIAKAD-Mini/app/model"
)

type EnrollmentRepository struct {
	Pool *pgxpool.Pool
}

var (
	ErrEnrollmentAlreadyExists = errors.New("enrollment already exists")
	ErrCourseNotFound          = errors.New("course not found")
	ErrCourseFull              = errors.New("course is full")
	ErrStudentNotFound         = errors.New("student not found")
	ErrSKSLimitExceeded        = errors.New("sks limit exceeded")
)

func (r *EnrollmentRepository) CreateEnrollment(
	ctx context.Context,
	studentID int,
	req model.CreateEnrollmentRequest,
) (*model.Enrollment, error) {

	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// 1. Cek mahasiswa dan tentukan batas SKS berdasarkan IPK
	var ipk float64

	err = tx.QueryRow(
		ctx,
		`SELECT COALESCE(ipk_terakhir, 0)
		 FROM students
		 WHERE id = $1
		   AND deleted_at IS NULL
		 FOR UPDATE`,
		studentID,
	).Scan(&ipk)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrStudentNotFound
		}
		return nil, err
	}

	batasSKS := 18

	if ipk >= 3.00 {
		batasSKS = 24
	} else if ipk >= 2.50 {
		batasSKS = 21
	}

	// 2. Cek apakah enrollment sudah ada
	var enrollment model.Enrollment

	err = tx.QueryRow(
		ctx,
		`SELECT id, student_id, course_id, tahun_akademik
		 FROM enrollments
		 WHERE student_id = $1
		   AND course_id = $2
		   AND tahun_akademik = $3`,
		studentID,
		req.CourseID,
		req.TahunAkademik,
	).Scan(
		&enrollment.ID,
		&enrollment.StudentID,
		&enrollment.CourseID,
		&enrollment.TahunAkademik,
	)

	if err == nil {
		return nil, ErrEnrollmentAlreadyExists
	}

	if err != pgx.ErrNoRows {
		return nil, err
	}

	// 3. Ambil data course dan kunci baris selama transaction
	var sks int
	var kuota int

	err = tx.QueryRow(
		ctx,
		`SELECT sks, kuota
		 FROM courses
		 WHERE id = $1
		 FOR UPDATE`,
		req.CourseID,
	).Scan(&sks, &kuota)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrCourseNotFound
		}
		return nil, err
	}

	// 4. Hitung jumlah mahasiswa yang sudah mengambil course
	var terisi int

	err = tx.QueryRow(
		ctx,
		`SELECT COUNT(*)
		 FROM enrollments
		 WHERE course_id = $1`,
		req.CourseID,
	).Scan(&terisi)

	if err != nil {
		return nil, err
	}

	if terisi >= kuota {
		return nil, ErrCourseFull
	}

	// 5. Hitung total SKS mahasiswa pada tahun akademik tersebut
	var totalSKS int

	err = tx.QueryRow(
		ctx,
		`SELECT COALESCE(SUM(c.sks), 0)
		 FROM enrollments e
		 JOIN courses c ON c.id = e.course_id
		 WHERE e.student_id = $1
		   AND e.tahun_akademik = $2`,
		studentID,
		req.TahunAkademik,
	).Scan(&totalSKS)

	if err != nil {
		return nil, err
	}

	// 6. Cek batas SKS
	// 6. Cek batas SKS
	if totalSKS+sks > batasSKS {
		sisaSKS := batasSKS - totalSKS

		return nil, fmt.Errorf(
			"%w: sisa SKS %d",
			ErrSKSLimitExceeded,
			sisaSKS,
		)
	}

	// 7. Simpan enrollment
	err = tx.QueryRow(
		ctx,
		`INSERT INTO enrollments
			(student_id, course_id, tahun_akademik)
		 VALUES ($1, $2, $3)
		 RETURNING id, student_id, course_id, tahun_akademik`,
		studentID,
		req.CourseID,
		req.TahunAkademik,
	).Scan(
		&enrollment.ID,
		&enrollment.StudentID,
		&enrollment.CourseID,
		&enrollment.TahunAkademik,
	)

	if err != nil {
		return nil, err
	}

	// 8. Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &enrollment, nil
}
