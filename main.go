package main

import (
	"log"
	"os"

	controller "github.com/NishantBansal2003/Brokerax/controller"
	routes "github.com/NishantBansal2003/Brokerax/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	PORT := os.Getenv("PORT")

	// Create a new Fiber app with the template engine
	app := fiber.New()
	// Set up CORS middleware
	app.Use(cors.New())
	// Connecting with MongoDB
	controller.Connect()

	// Setup the routes
	routes.Setup(app)

	// Set up flash middleware
	app.Get("/api", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"success": true,
			"message": "Hello World",
		})
	})
	log.Fatal(app.Listen(":" + PORT))
}
