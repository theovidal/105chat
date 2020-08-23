package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

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
	log.Println("Starting 105chat...")

	go ws.HandlePipeline()
	http.Handle("/v1/ws", websocket.Handler(ws.Server))

	go httpServer.Server()

	log.Println("WebSocket server ready")
	err := http.ListenAndServe("localhost:1051", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func Migrate(_ []string) {
	log.Println("Database migration started")
	db.Database.AutoMigrate(&db.User{}, &db.Message{}, &db.Room{})
	log.Println("Database migration complete")
}
