package model

type CreateEnrollmentRequest struct {
	CourseID      int    `json:"course_id"`
	TahunAkademik string `json:"tahun_akademik"`
}

type Enrollment struct {
	ID            int    `json:"id"`
	StudentID     int    `json:"student_id"`
	CourseID      int    `json:"course_id"`
	TahunAkademik string `json:"tahun_akademik"`
}
