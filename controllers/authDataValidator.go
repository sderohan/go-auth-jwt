package controllers

import (
	"errors"
	"strings"
	"unicode"
)

var (
	ErrInvalidNameFormat     = errors.New("name does not match the specified constraints")
	ErrInvalidEmailFormat    = errors.New("email format is invalid")
	ErrInvalidPasswordFormat = errors.New("password does not match specified the constraints")

	ErrInvalidEmail    = errors.New("user with given email does not exist")
	ErrInvalidPassword = errors.New("incorrect password")
	ErrUnauthenticated = errors.New("unauthenticated")

	ErrInternal = errors.New("error occured, please retry")
)

type valData func(data *map[string]string) error

func validateData(data *map[string]string, fns ...valData) error {
	for _, fn := range fns {
		if err := fn(data); err != nil {
			return err
		}
	}
	return nil
}

/*
 Below function validates the username from the request body
 only uppercase and lowercase letters are allowed
 and length of the name should be greater than 4 (could be less or more depending on the usecase)
*/
func validateName(data *map[string]string) error {
	if !isExist("name", data) {
		return ErrInvalidNameFormat
	}
	name := (*data)["name"]
	name_len := len(name)
	if name_len < 5 {
		return ErrInvalidNameFormat
	}
	for _, char := range name {
		if !unicode.IsLetter(char) {
			return ErrInvalidNameFormat
		}
	}
	return nil
}

func validateEmail(data *map[string]string) error {
	if !isExist("email", data) {
		return ErrInvalidEmailFormat
	}
	// Add below the email validation code
	return nil
}

/*
	This function validate the user password and returns error
	if password does not match the specified contraints

	Following are the constraints the password string should meet
	The password string should have atleast
	- length 8
	- one uppercase letter
	- one lowercase letter
	- one digit
	- one special symbol specified in the ALLOWED_SPECIAL_CHARS variable inside function
*/
func validatePassword(data *map[string]string) error {

	// first check if password field is set into the map
	// this is to check if user has sent the password field in the json request
	if !isExist("password", data) {
		return ErrInvalidPasswordFormat
	}

	// password string to process
	str := (*data)["password"]

	// remove leading and trailing whitespaces if any
	str = strings.Trim(str, " ")

	// get the length of the processed string
	n := len(str)

	// constraints
	var (
		UPPERCASE    = "UPPERCASE"
		LOWERCASE    = "LOWERCASE"
		DIGIT        = "DIGIT"
		SPECIAL_CHAR = "SPECIAL_CHAR"
	)

	// map of constraints, initially set to false
	constraints := map[string]bool{
		UPPERCASE:    false,
		LOWERCASE:    false,
		DIGIT:        false,
		SPECIAL_CHAR: false,
	}

	// number of constraints
	number_of_contraints := 4

	// list of allowed special characters in password
	ALLOWED_SPECIAL_CHARS := "!@#$%^&"

	// minimum password length
	password_length := 8

	// string length should be >= 8
	if n >= password_length {
		// check each char and update the constraint
		for _, char := range str {
			// Check each constraint, if not set then set it and update the count of number of constraint
			// This is done so that we don't have to iterate the entire string chars if all constraints are already set
			switch {
			case unicode.IsDigit(char) && !isConstraintSet(constraints, DIGIT):
				constraints[DIGIT] = true
				number_of_contraints--
			case unicode.IsUpper(char) && !isConstraintSet(constraints, UPPERCASE):
				constraints[UPPERCASE] = true
				number_of_contraints--
			case unicode.IsLower(char) && !isConstraintSet(constraints, LOWERCASE):
				constraints[LOWERCASE] = true
				number_of_contraints--
			case strings.Contains(ALLOWED_SPECIAL_CHARS, string(char)) && !isConstraintSet(constraints, SPECIAL_CHAR):
				constraints[SPECIAL_CHAR] = true
				number_of_contraints--
			}
			// check if all constraints are set and break
			if number_of_contraints == 0 {
				break
			}
		}
		// iterate over the constraints and check if all constraints are set to true
		for _, value := range constraints {
			// if any one of the constraint is not set then return error
			if !value {
				return ErrInvalidPasswordFormat
			}
		}
		// error is nil when all constraints are set to true
		return nil
	}
	return ErrInvalidPasswordFormat
}

// Check if the given key exist in the map
func isExist(data string, mp *map[string]string) bool {
	_, exist := (*mp)[data]
	return exist
}

// returns the value of the key in map
func isConstraintSet(constraints map[string]bool, name string) bool {
	return constraints[name]
}
