package main

import (
	"log"
	"net/http"

	"github.com/fatih/color"
	"golang.org/x/net/websocket"

	"github.com/theovidal/105chat/cache"
	"github.com/theovidal/105chat/db"
	httpServer "github.com/theovidal/105chat/http"
	"github.com/theovidal/105chat/utils"
	"github.com/theovidal/105chat/ws"
)

// Run is a CLI command that starts the 105chat API
func Run(_ []string) {
	color.HiCyan("\n  _____   ________   ________   ________   ___  ___   ________   _________   \n / __  \\ |\\   __  \\ |\\   ____\\ |\\   ____\\ |\\  \\|\\  \\ |\\   __  \\ |\\___   ___\\ \n|\\/_|\\  \\\\ \\  \\|\\  \\\\ \\  \\___|_\\ \\  \\___| \\ \\  \\\\\\  \\\\ \\  \\|\\  \\\\|___ \\  \\_| \n\\|/ \\ \\  \\\\ \\  \\\\\\  \\\\ \\_____  \\\\ \\  \\     \\ \\   __  \\\\ \\   __  \\    \\ \\  \\  \n     \\ \\  \\\\ \\  \\\\\\  \\\\|____|\\  \\\\ \\  \\____ \\ \\  \\ \\  \\\\ \\  \\ \\  \\    \\ \\  \\ \n      \\ \\__\\\\ \\_______\\ ____\\_\\  \\\\ \\_______\\\\ \\__\\ \\__\\\\ \\__\\ \\__\\    \\ \\__\\\n       \\|__| \\|_______||\\_________\\\\|_______| \\|__|\\|__| \\|__|\\|__|     \\|__|\n                       \\|_________|                                          \n")

	log.Println(color.CyanString("⏩ Step 1: Open databases"))
	db.OpenDatabase()
	cache.OpenCache()

	log.Println(color.CyanString("⏩ Step 2: Populate cache"))
	Cache(nil)

	log.Println(color.CyanString("⏩ Step 3: Start HTTP server"))
	go httpServer.Server()

	log.Println(color.CyanString("⏩ Step 4: Start WS server"))
	wsServer := ws.NewServer()
	go wsServer.Listen()
	http.Handle("/v1/ws", websocket.Handler(wsServer.Handle))

	addr := utils.GenerateAddress("WS")
	log.Println("▶ WS server listening on", color.CyanString("ws://"+addr))
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Panicf("‼ WS server fatal error: %s", err.Error())
	}
}
