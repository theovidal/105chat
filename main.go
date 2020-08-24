package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"golang.org/x/net/websocket"

	"github.com/theovidal/105chat/db"
	httpServer "github.com/theovidal/105chat/http"
	"github.com/theovidal/105chat/ws"
)

type Command struct {
	Name        string
	Description string
	FlagSet     *flag.FlagSet
	Handler     func([]string)
}

func (c *Command) String() string {
	return fmt.Sprintf("%s - %s", c.Name, c.Description)
}

var commands = make(map[string]Command)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file at the root - Ignoring")
	}

	commands = map[string]Command{
		"help": Command{
			Name:        "help",
			Description: "Show help",
			FlagSet:     flag.NewFlagSet("help", flag.ExitOnError),
			Handler:     Help,
		},
		"run": Command{
			Name:        "run",
			Description: "Run the 105chat server",
			FlagSet:     flag.NewFlagSet("run", flag.ExitOnError),
			Handler:     Run,
		},
		"migrate": Command{
			Name:        "migrate",
			Description: "Migrate database models",
			FlagSet:     flag.NewFlagSet("migrate", flag.ExitOnError),
			Handler:     Migrate,
		},
	}

	if len(os.Args) < 2 {
		Help([]string{})
		os.Exit(0)
	}

	command, found := commands[os.Args[1]]
	if !found {
		fmt.Println("Command", os.Args[1], "is not valid. Run help command to get the full list of commands")
		os.Exit(1)
	}

	command.Handler(command.FlagSet.Args())
}

func Help(_ []string) {
	fmt.Println("───── 105chat CLI help ─────\n")

	fmt.Println("COMMANDS")
	for _, command := range commands {
		fmt.Println(command.String())
	}
}

func Run(_ []string) {
	color.HiGreen("\n  _____   ________   ________   ________   ___  ___   ________   _________   \n / __  \\ |\\   __  \\ |\\   ____\\ |\\   ____\\ |\\  \\|\\  \\ |\\   __  \\ |\\___   ___\\ \n|\\/_|\\  \\\\ \\  \\|\\  \\\\ \\  \\___|_\\ \\  \\___| \\ \\  \\\\\\  \\\\ \\  \\|\\  \\\\|___ \\  \\_| \n\\|/ \\ \\  \\\\ \\  \\\\\\  \\\\ \\_____  \\\\ \\  \\     \\ \\   __  \\\\ \\   __  \\    \\ \\  \\  \n     \\ \\  \\\\ \\  \\\\\\  \\\\|____|\\  \\\\ \\  \\____ \\ \\  \\ \\  \\\\ \\  \\ \\  \\    \\ \\  \\ \n      \\ \\__\\\\ \\_______\\ ____\\_\\  \\\\ \\_______\\\\ \\__\\ \\__\\\\ \\__\\ \\__\\    \\ \\__\\\n       \\|__| \\|_______||\\_________\\\\|_______| \\|__|\\|__| \\|__|\\|__|     \\|__|\n                       \\|_________|                                          \n")
	db.OpenDatabase()

	go ws.HandlePipeline()
	http.Handle("/v1/ws", websocket.Handler(ws.Server))

	go httpServer.Server()

	addr := os.Getenv("WS_ADDRESS") + ":" + os.Getenv("WS_PORT")
	log.Println("WS server listening on", color.CyanString(addr))
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Panicf("WS server fatal error: %s", err.Error())
	}
}

func Migrate(_ []string) {
	log.Println("Database migration started")
	db.Database.AutoMigrate(&db.User{}, &db.Message{}, &db.Room{})
	log.Println("Database migration complete")
}
