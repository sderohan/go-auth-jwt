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

// informational messages
var (
	LOGIN_SUCCESSFULL  = "logged in successfully"
	LOGOUT_SUCCESSFULL = "logged out successfully"
	USER_EXIST         = "user account already exists with the same email id!"
)

// Secret is used to sign the jwt token, this is dummy secret.
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

	// if data is invalid, return the error
	if err != nil {
		return sendResponse(c, err.Error())
	}

	// Check if user exists in the database with given email id
	var user models.User
	userExist("email = ?", data["email"], &user)
	if user.Email == data["email"] {
		return sendResponse(c, USER_EXIST)
	}

	// Hash the password
	password, err := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	if err != nil {
		return err
	}

	// Create the user object
	user = models.User{
		Name:     data["name"],
		Email:    data["email"],
		Password: string(password),
	}

	// Insert the user in the database
	database.DB.Create(&user)

	// Return the newly created user information in the database
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

	// If data is in invalid format, send the error message
	if err != nil {
		return sendResponse(c, err.Error())
	}

	var user models.User

	// Load the user with the given email matching from the database
	getFirst("email = ?", data["email"], &user)

	// look for the user with the given email in the database if exists
	if user.ID == 0 {
		return sendResponse(c, ErrInvalidEmail.Error(), fiber.StatusNotFound)
	}

	// Check if the given password match with the password in the database
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"])); err != nil {
		return sendResponse(c, ErrInvalidPassword.Error(), fiber.StatusBadRequest)
	}

	// cookie expiration time
	cookie_expiry := time.Now().Add(time.Hour * 24)

	// Create the jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": json.Number(strconv.FormatInt(cookie_expiry.Unix(), 10)),
		"iss": fmt.Sprintf("%d", user.ID),
	})

	// Sign the newly created token with the secret
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return sendResponse(c, "Error occured, please retry", fiber.StatusInternalServerError)
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

	return sendResponse(c, LOGIN_SUCCESSFULL, fiber.StatusOK)
}

// Function returns the user information
// GET
func User(c *fiber.Ctx) error {

	// Get the cookie
	cookie := c.Cookies("jwt")

	// Validate the user and get the token
	token, err := validateToken(c, cookie)
	if err != nil {
		return sendResponse(c, err.Error(), fiber.StatusUnauthorized)
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User

	// load the user with matching user ID
	getFirst("id = ?", claims.Issuer, &user)

	// return the user information
	return c.JSON(user)
}

// Logout clear the user
func Logout(c *fiber.Ctx) error {
	// Set the new cookie, altering the time
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return sendResponse(c, LOGOUT_SUCCESSFULL)
}
