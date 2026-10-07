package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/ar-elfahmi/SIAKAD-Mini/app/model"
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

	studentRepository := repository.StudentRepository{
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
	app.Get("/api/v1/students", func(c *fiber.Ctx) error {
		page := c.QueryInt("page", 1)
		if page < 1 {
			page = 1
		}
		perPage := c.QueryInt("per_page", 10)
		if perPage < 1 {
			perPage = 10
		}
		if perPage > 50 {
			perPage = 50
		}
		prodi := c.Query("prodi")
		angkatan := c.QueryInt("angkatan", 0)
		search := c.Query("search")
		sort := c.Query("sort")

		students, total, err := studentRepository.GetAllStudents(
			c.Context(),
			page,
			perPage,
			prodi,
			angkatan,
			search,
			sort,
		)
		if err != nil {
			return err
		}

		lastPage := (total + perPage - 1) / perPage

		return c.JSON(fiber.Map{
			"data": students,
			"meta": model.StudentMeta{
				CurrentPage: page,
				PerPage:     perPage,
				Total:       total,
				LastPage:    lastPage,
			},
		})
	})
	app.Get("/api/v1/students/:id", func(c *fiber.Ctx) error {
		id, parseErr := c.ParamsInt("id")
		if parseErr != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "student not found",
			})
		}
		student, err := studentRepository.GetStudentByID(c.Context(), id)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "student not found",
				})
			}
			return err
		}
		return c.JSON(student)
	})
	app.Listen(":3000")
}
