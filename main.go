package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/iqbaleff214/easynote-backend-go/auth"
	"github.com/iqbaleff214/easynote-backend-go/handler"
	"github.com/iqbaleff214/easynote-backend-go/user"
)

func main() {
	initConfig()

	// database init
	db, err := database(appConfig.mysqlUri)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { db.Close() }()

	// repository init
	userRepository := user.NewRepository(db)

	// service init
	authService := auth.NewService(appConfig.jwtSecret)
	userService := user.NewService(userRepository)

	// handler init
	userHandler := handler.NewUserHandler(userService, authService)

	app := fiber.New()
	app.Use(cors.New())

	api := app.Group("/api/v1", logger.New())

	// Public
	api.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(map[string]string{
			"Product Name": "EasyNote",
			"Version":      appConfig.version,
			"Date":         "30/01/2024",
		})
	})

	// User Domain
	api.Post("/register", userHandler.RegisterUser)
	api.Post("/login", userHandler.Login)

	api.Use(handler.AuthMiddleware(appConfig.jwtSecret, userService))

	api.Get("/profile", userHandler.CurrentUser)
	api.Put("/profile", userHandler.UpdateUser)

	// Note Domain

	// Folder Domain

	log.Fatal(app.Listen(":8000"))
}
