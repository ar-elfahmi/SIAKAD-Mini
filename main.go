package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env")
	}
		
	databaseURL := os.Getenv("DATABASE_URL")

	poll, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	defer poll.Close()

	err = poll.Ping(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	
	log.Println("Database Connected")

	courses := []Course{}
	
	rows, err := poll.Query(
		context.Background(),
		`SELECT id, kode_mk, nama_mk, sks, semester, kuota FROM courses`,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var course Course

		err := rows.Scan(
			&course.ID,
			&course.KodeMK,
			&course.NamaMK,
			&course.SKS,
			&course.Semester,
			&course.Kuota,
		)

		if err != nil {
			log.Fatal(err)
		}

		courses = append(courses, course)
	}

	if err := rows.Err(); err != nil{
		log.Fatal(err)
	}
	
	app := fiber.New()

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
