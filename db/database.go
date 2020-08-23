package db

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/qor/validations"
)

// Database is the database used to store static content (messages, rooms, users...).
// For now it's a SQLite DB, change to MySQL or PostgreSQL is planned
var Database = openDatabase()

// openDatabase setups the database, and create the file if it doesn't exists
func openDatabase() *gorm.DB {
	filename := "testing.sqlite"

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		_, err := os.Create(filename)
		if err != nil {
			log.Panic(fmt.Sprintf("Error while creating the database : %s", err))
		}
	}

	db, err := gorm.Open("sqlite3", filename)
	if err != nil {
		log.Panic(fmt.Sprintf("failed to connect database : %s", err))
	}
	validations.RegisterCallbacks(db)

	return db
}
