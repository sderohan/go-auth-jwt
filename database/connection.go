package database

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	host          = "localhost"
	port          = 49154
	user_name     = "root"
	password      = "root"
	database_name = "go_auth"
)

func Connect() {

	// Database connection string
	connectionString := fmt.Sprintf("%s:%s@(%s:%d)/%s?parseTime=true", user_name, password, host, port, database_name)

	// Connect to the database container
	_, err := gorm.Open(mysql.Open(connectionString))
	if err != nil {
		log.Fatal("Could not create the database connection ", err.Error())
	}
}
