package main

import (
	"log"

	"github.com/fatih/color"

	"github.com/theovidal/105chat/db"
)

// Migrate is a CLI command that performs a migration on the database to register models
func Migrate(_ []string) {
	db.OpenDatabase()
	log.Println("✏ Database migration started")
	db.Client.AutoMigrate(
		&db.User{},
		&db.Message{},
		&db.Room{},
		&db.Group{},
		&db.RoomPermission{},
		&db.GroupInheritance{},
	)
	log.Println(color.GreenString("✅ Database migration complete"))
}
