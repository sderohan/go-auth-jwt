package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sderohan/go-auth-jwt/database"
	"github.com/sderohan/go-auth-jwt/models"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	password, err := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	if err != nil {
		return err
	}

	user := models.User{
		Name:     data["name"],
		Email:    data["email"],
		Password: string(password),
	}

	// Create the user in the database
	database.DB.Create(&user)

	return c.JSON(user)
}
