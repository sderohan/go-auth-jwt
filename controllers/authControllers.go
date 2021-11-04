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

	// Parse the request
	parseData(c, &data)

	// Validate the data
	err := validateData(
		&data,
		validateName,
		validateEmail,
		validatePassword,
	)

	if err != nil {
		return sendResponse(c, err.Error())
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

	// Parse the request
	parseData(c, &data)

	// Validate the data
	err := validateData(
		&data,
		validateEmail,
		validatePassword,
	)

	if err != nil {
		return sendResponse(c, err.Error())
	}

	var user models.User

	// Load the user with the given email matching from the database
	getFirst("email = ?", data["email"], &user)

	// look for the user with the given email in the database if exists
	if user.ID == 0 {
		return sendResponse(c, "User with given email does not exist", fiber.StatusNotFound)
	}

	// check if the given password match with the password in the database
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"])); err != nil {
		return sendResponse(c, "incorrect password", fiber.StatusBadRequest)
	}

	// cookie expiration time
	cookie_expiry := time.Now().Add(time.Hour * 24)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": json.Number(strconv.FormatInt(cookie_expiry.Unix(), 10)),
		"iss": fmt.Sprintf("%d", user.ID),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return sendResponse(c, "Could not login", fiber.StatusInternalServerError)
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

	return sendResponse(c, "success", fiber.StatusOK)
}

// Function returns the user information
// GET
func User(c *fiber.Ctx) error {

	// Get the cookie
	cookie := c.Cookies("jwt")

	// Validate the user and get the token
	token, _ := validateToken(c, cookie)

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User

	// load the user with matching user ID
	getFirst("id = ?", claims.Issuer, &user)

	// return the user information
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
	return sendResponse(c, "Logged out successfully")
}
