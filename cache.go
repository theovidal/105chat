package main

import (
	"github.com/fatih/color"
	"log"

	"github.com/theovidal/105chat/db"
)

func Cache(_ []string) {
	log.Println("📁 Cache population started")
	db.SetAllGroupsCache()
	log.Println(color.HiGreenString("✅ Cache population completed"))
}
