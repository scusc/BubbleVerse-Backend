package db

import (
	"fmt"
	"log"
	"os"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/joho/godotenv"
)

// DB is a global variable that holds a pointer to a gorm.DB instance.
// gorm.DB is the main database handle provided by the GORM library, which is a popular ORM (Object Relational Mapper) for Go.
// This handle is used to interact with your database (e.g., running queries, inserting, updating, deleting records).
// By declaring it as a global variable, you can access the database connection from anywhere in your project.
// Typical usage:
// 1. Initialize DB (usually in your main function or an init function) by connecting to your database.
// 2. Use DB to perform database operations throughout your application.
// Example usage elsewhere:
//   DB.Find(&users) --> Retrieves all users from the database
var DB *gorm.DB

func InitDB() {
	godotenv.Load()
	// data source name (DSN) for PostgreSQL
	// This would typically come from environment variables or a configuration file
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Los_Angeles",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	log.Println("Connected to DB")
}
