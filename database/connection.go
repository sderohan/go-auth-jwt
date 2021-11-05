package database

import (
	"fmt"
	"log"

	"github.com/sderohan/go-auth-jwt/models"
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

// Shared across package to perform the database operations
var DB *gorm.DB

func Connect() {

	// Database connection string
	connectionString := fmt.Sprintf("%s:%s@(%s:%d)/%s?parseTime=true", user_name, password, host, port, database_name)

	// Connect to the database container
	connection, err := gorm.Open(mysql.Open(connectionString))

	// Alternatively err variable can be declared first but I preffered this
	DB = connection
	if err != nil {
		log.Fatal("Could not create the database connection ", err.Error())
	}

	// Migrate the user in database
	DB.AutoMigrate(&models.User{})
}
