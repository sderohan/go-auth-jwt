package database

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	HOST          = "localhost"
	PORT          = 49154
	USER_NAME     = "root"
	PASSWORD      = "root"
	DATABASE_NAME = "go_auth"
)

func Connect() {

	// Database connection string
	connectionString := fmt.Sprintf("%s:%s@(%s:%d)/%s?parseTime=true", USER_NAME, PASSWORD, HOST, PORT, DATABASE_NAME)

	// Connect to the database container
	_, err := gorm.Open(mysql.Open(connectionString))
	if err != nil {
		log.Fatal("Could not create the database connection ", err.Error())
	}
}
