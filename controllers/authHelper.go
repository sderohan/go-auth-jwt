package controllers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/sderohan/go-auth-jwt/database"
	"github.com/sderohan/go-auth-jwt/models"
)

// parseData() processes the request data and returns it to the destination map
func parseData(c *fiber.Ctx, dst *map[string]string) error {
	if err := c.BodyParser(&dst); err != nil {
		return err
	}
	return nil
}

// Returns the first matching user information from the database
func getFirst(query string, data interface{}, user *models.User) {
	database.DB.Where(query, data).First(user)
}

// Check the uer exist in the database with the given email id
func userExist(query string, data interface{}, user *models.User) {
	database.DB.Where(query, data).Find(user)
}

// Checks if the current user token is valid and returns it if valid
func validateToken(c *fiber.Ctx, cookie string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, ErrUnauthenticated
	}
	return token, nil
}

func sendResponse(c *fiber.Ctx, msg string, statusCode ...int) error {
	if len(statusCode) > 0 {
		c.Status(statusCode[0])
	}
	return c.JSON(fiber.Map{
		"message": msg,
	})
}
