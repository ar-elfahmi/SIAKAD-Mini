package model

type Student struct {
	ID          int     `json:"id"`
	NIM         string  `json:"nim"`
	Nama        string  `json:"nama"`
	Prodi       string  `json:"prodi"`
	Angkatan    int     `json:"angkatan"`
	IPKTerakhir float64 `json:"ipk_terakhir"`
}

type CreateStudentRequest struct {
	NIM         string   `josn:"nim"`
	Nama        string   `josn:"nama"`
	Email       string   `josn:"email"`
	Prodi       string   `josn:"prodi"`
	Angkatan    int      `josn:"angkatan"`
	IPKTerakhir *float64 `josn:"ipk_terakhir"`
}

type StudentMeta struct {
	CurrentPage int `json:"current_page"`
	PerPage     int `json:"per_page"`
	Total       int `json:"total"`
	LastPage    int `json:"last_page"`
}

type EnrolledCourse struct {
	ID            int    `json:"id"`
	KodeMK        string `json:"kode_mk"`
	NamaMK        string `json:"nama_mk"`
	SKS           int    `json:"sks"`
	Semester      int    `json:"semester"`
	TahunAkademik string `json:"tahun_akademik"`
}

type StudentDetail struct {
	ID          int              `json:"id"`
	NIM         string           `json:"nim"`
	Nama        string           `json:"nama"`
	Prodi       string           `json:"prodi"`
	Angkatan    int              `json:"angkatan"`
	IPKTerakhir float64          `json:"ipk_terakhir"`
	TotalSKS    int              `json:"total_sks"`
	BatasSKS    int              `json:"batas_sks"`
	Courses     []EnrolledCourse `json:"courses"`
}
