package controllers

import "errors"

var (
	ErrInvalidNameFormat     = errors.New("name does not match the constraints")
	ErrInvalidEmailFormat    = errors.New("email format is not valid")
	ErrInvalidPasswordFormat = errors.New("password does not match the constraints")

	ErrInvalidEmail    = errors.New("user with given email does not exist")
	ErrInvalidPassword = errors.New("incorrect password")
	ErrUnauthenticated = errors.New("unauthenticated")
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

func validateName(data *map[string]string) error {
	if !isExist("name", data) {
		return ErrInvalidNameFormat
	}
	// Need to add the name validation
	return nil
}

func validateEmail(data *map[string]string) error {
	if !isExist("email", data) {
		return ErrInvalidEmailFormat
	}
	// Add below the email validation code
	return nil
}

func validatePassword(data *map[string]string) error {
	if !isExist("password", data) {
		return ErrInvalidPasswordFormat
	}
	// Add below the password validation code
	return nil
}

func isExist(data string, mp *map[string]string) bool {
	_, exist := (*mp)[data]
	return exist
}
