package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/ar-elfahmi/SIAKAD-Mini/app/repository"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env")
	}

	databaseURL := os.Getenv("DATABASE_URL")

	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	defer pool.Close()

	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database Connected")

	courseRepository := repository.CourseRepository{
		Pool: pool,
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("yey my first hello world")
	})
	app.Get("/api/v1/courses", func(c *fiber.Ctx) error {
		semester := c.QueryInt("semester", 0)
		search := c.Query("search")
		available := c.QueryBool("available", false)
		
		courses, err := courseRepository.GetAllCourses(c.Context(), semester, 
		search,
		available,)
		if err != nil {
			return err
		}
		return c.JSON(courses)
	})
	app.Listen(":3000")
}
