package db

import (
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/qor/validations"
)

// Database is the database used to store static content (messages, rooms, users...).
// For now it's a SQLite DB, change to MySQL or PostgreSQL is planned
var Database *gorm.DB

// OpenDatabase setups the database, and create the file if it doesn't exists
func OpenDatabase() {
	filename := os.Getenv("DB_FILE")

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		_, err := os.Create(filename)
		if err != nil {
			log.Fatalf("Error while creating the database: %s", err)
		}
	}

	var err error
	Database, err = gorm.Open("sqlite3", filename)
	if err != nil {
		log.Fatalf("Failed to open database: %s", err)
	}
	Database.LogMode(os.Getenv("DEBUG") == "on")
	validations.RegisterCallbacks(Database)
}
