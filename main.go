package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sderohan/go-auth-jwt/database"
	"github.com/sderohan/go-auth-jwt/routes"
)

func main() {

	// Connect to the database
	database.Connect()

	// Create the app instance
	app := fiber.New()

	// Enable the CORS so that the frontend running on other port can access it
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	// Register the routes
	routes.Setup(app)

	// Start the server
	app.Listen(":3000")
}
