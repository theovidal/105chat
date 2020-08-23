package db

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/qor/validations"
)

var Database = openDb()

func openDb() *gorm.DB {
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
