package model

type Student struct {
	ID         int     `json:"id"`
	NIM        string  `json:"nim"`
	Nama       string  `json:"nama"`
	Prodi      string  `json:"prodi"`
	Angkatan   int     `json:"angkatan"`
	IPKTerakhir float64 `json:"ipk_terakhir"`
}

type StudentMeta struct {
	CurrentPage int `json:"current_page"`
	PerPage     int `json:"per_page"`
	Total       int `json:"total"`
	LastPage    int `json:"last_page"`
}
