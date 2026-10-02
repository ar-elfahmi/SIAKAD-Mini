package main

import "github.com/gofiber/fiber/v2"

func main() {
	app := fiber.New()

	courses := []Course{
		{
			ID:       1,
			KodeMK:   "IF101",
			NamaMK:   "Pemerograman Backend Lanjut Praktikum",
			SKS:      "3",
			Semester: "5",
			Kuota:    "40",
		},
		{
			ID:       2,
			KodeMK:   "ML032",
			NamaMK:   "Mesin Learning",
			SKS:      "2",
			Semester: "5",
			Kuota:    "40",
		},
		{
			ID:       1,
			KodeMK:   "IF100",
			NamaMK:   "Pemerograman Backend Lanjut Teori",
			SKS:      "1",
			Semester: "5",
			Kuota:    "40",
		},
		{
			ID:       1,
			KodeMK:   "DT001",
			NamaMK:   "Design Thingking",
			SKS:      "2",
			Semester: "5",
			Kuota:    "40",
		},
		{
			ID:       1,
			KodeMK:   "KWU13",
			NamaMK:   "Kewirausahaan",
			SKS:      "2",
			Semester: "5",
			Kuota:    "40",
		},
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("yey my first hello world")
	})
	app.Get("/api/v1/courses", func(c *fiber.Ctx) error {
		return c.JSON(courses)
	})
	app.Listen(":3000")
}

type Course struct {
	ID       int    `json:"id"`
	KodeMK   string `json:"kode_mk"`
	NamaMK   string `json:"nama_mk"`
	SKS      string `json:"sks"`
	Semester string `json:"semester"`
	Kuota    string `json:"kuota"`
}
