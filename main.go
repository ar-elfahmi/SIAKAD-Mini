package main

import (
	"context"
	"errors"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	"github.com/ar-elfahmi/SIAKAD-Mini/app/middleware"
	"github.com/ar-elfahmi/SIAKAD-Mini/app/model"
	"github.com/ar-elfahmi/SIAKAD-Mini/app/repository"
	"github.com/ar-elfahmi/SIAKAD-Mini/app/service"
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

	authRepository := repository.AuthRepository{
		Pool: pool,
	}

	authService := service.AuthService{
		JWTSecret: os.Getenv("JWT_SECRET"),
	}

	app := fiber.New()

	authMiddleware := middleware.RequireAuth(
		os.Getenv("JWT_SECRET"),
	)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("yey my first hello world")
	})
	app.Post("/api/v1/auth/login", func(c *fiber.Ctx) error {
		var req model.LoginRequest

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error": "invalid request body",
			})
		}

		user, passwordHash, err := authRepository.FindUserByEmail(
			c.Context(),
			req.Email,
		)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "email atau password salah",
			})
		}

		if err := bcrypt.CompareHashAndPassword(
			[]byte(passwordHash),
			[]byte(req.Password),
		); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "email atau password salah",
			})
		}

		token, expiresIn, err := authService.GenerateToken(user)
		if err != nil {
			return err
		}

		return c.JSON(model.LoginResponse{
			AccessToken: token,
			TokenType:   "Bearer",
			ExpiresIn:   expiresIn,
			User:        *user,
		})
	})
	app.Get("/api/v1/auth/me", authMiddleware, func(c *fiber.Ctx) error {
		userIDFloat, ok := c.Locals("user_id").(float64)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "token tidak valid",
			})
		}

		userID := int(userIDFloat)

		user, err := authRepository.GetUserByID(
			c.Context(),
			userID,
		)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "user tidak ditemukan",
			})
		}

		return c.JSON(user)
	})
	app.Get("/api/v1/courses", authMiddleware, func(c *fiber.Ctx) error {
		semester := c.QueryInt("semester", 0)
		search := c.Query("search")
		available := c.QueryBool("available", false)

		courses, err := courseRepository.GetAllCourses(c.Context(), semester,
			search,
			available)
		if err != nil {
			return err
		}
		return c.JSON(courses)
	})
	app.Get(
		"/api/v1/students",
		authMiddleware,
		middleware.RequireRole("admin"),
		func(c *fiber.Ctx) error {
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
	app.Get("/api/v1/students/:id", authMiddleware, func(c *fiber.Ctx) error {
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
	app.Post("/api/v1/students", authMiddleware, func(c *fiber.Ctx) error {
		var req model.CreateStudentRequest

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error": "invalid request body",
			})
		}

		if !regexp.MustCompile(`^\d{12}$`).MatchString(req.NIM) {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error": "nim harus 12 digit",
			})
		}
		if strings.TrimSpace(req.Nama) == "" {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error": "nama wajib diisi",
			})
		}

		if strings.TrimSpace(req.Email) == "" || !strings.Contains(req.Email, "@") {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error": "email tidak valid",
			})
		}

		if strings.TrimSpace(req.Prodi) == "" {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error": "prodi wajib diisi",
			})
		}

		currentYear := time.Now().Year()

		if req.Angkatan < 1000 || req.Angkatan > currentYear {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error": "angkatan tidak valid",
			})
		}

		if req.IPKTerakhir != nil &&
			(*req.IPKTerakhir < 0 || *req.IPKTerakhir > 4) {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error": "ipk_terakhir harus antara 0 sampai 4",
			})
		}

		passwordHash, err := bcrypt.GenerateFromPassword(
			[]byte(req.NIM),
			bcrypt.DefaultCost,
		)
		if err != nil {
			return err
		}

		student, err := studentRepository.CreateStudent(
			c.Context(),
			req,
			string(passwordHash),
		)

		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error": "nim atau email sudah terdaftar",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(student)
	})
	app.Listen(":3000")
}
