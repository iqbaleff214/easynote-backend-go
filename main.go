package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func main() {
	db, err := database("root:@tcp(127.0.0.1:3306)/easynote?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer func ()  { db.Close() }()

	app := fiber.New()
	app.Use(cors.New())

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello World")
	})

	log.Fatal(app.Listen(":8000"))
}