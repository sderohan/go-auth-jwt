package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sderohan/go-auth-jwt/database"
	"github.com/sderohan/go-auth-jwt/routes"
)

func main() {

	// Connect to the database
	database.Connect()

	// Create the app instance
	app := fiber.New()

	// Register the routes
	routes.Setup(app)

	// Start the server
	app.Listen(":3000")
}
