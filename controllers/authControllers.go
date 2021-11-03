package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/sderohan/go-auth-jwt/database"
	"github.com/sderohan/go-auth-jwt/models"
	"golang.org/x/crypto/bcrypt"
)

// secret is used to sign the jwt token, this is dummy secret.
// Change it as per your need
const secret = "strongsecretmessage"

// This function register's the new user
// POST
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

// Allows the user to login into the system using the existing credentials
// POST
func Login(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User

	// Load the user with the given email matching from the database
	database.DB.Where("email = ?", data["email"]).First(&user)

	// look for the user with the given email in the database if exists
	if user.ID == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "User with given email does not exist",
		})
	}

	// check if the given password match with the password in the database
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "incorrect password",
		})
	}

	// cookie expiration time
	cookie_expiry := time.Now().Add(time.Hour * 24)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": json.Number(strconv.FormatInt(cookie_expiry.Unix(), 10)),
		"iss": fmt.Sprintf("%d", user.ID),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println(err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Couldn't login",
		})
	}

	// Create the cookie
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Expires:  cookie_expiry,
		HTTPOnly: true, // true, doesn't allow the web browser to modify the cookie
	}

	// set the cookie
	c.Cookie(&cookie)

	c.Status(fiber.StatusOK)
	return c.JSON(fiber.Map{
		"message": "success",
	})
}

// Function returns the user information
// GET
func User(c *fiber.Ctx) error {

	// Get the cookie
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User

	database.DB.Where("id = ?", claims.Issuer).First(&user)

	return c.JSON(user)
}

// Logout clear the user
func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return c.JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}
