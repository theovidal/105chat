package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"golang.org/x/net/websocket"

	"github.com/theovidal/105chat/db"
	httpServer "github.com/theovidal/105chat/http"
	"github.com/theovidal/105chat/ws"
)

// Command is the holder for all CLI commands
type Command struct {
	// The name of the command
	Name string
	// A short and efficient description for the command
	Description string
	// The flag set to register with the command
	FlagSet *flag.FlagSet
	// The handler function when the command is triggered
	Handler func([]string)
}

// String gives a string representation of the command
// Used to get help
func (c *Command) String() string {
	return fmt.Sprintf("%s - %s", c.Name, c.Description)
}

// commands holds all of the CLI commands
var commands = make(map[string]Command)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("💾 No .env file at the root - Ignoring")
	}

	commands = map[string]Command{
		"help": {
			Name:        "help",
			Description: "Show help",
			FlagSet:     flag.NewFlagSet("help", flag.ExitOnError),
			Handler:     Help,
		},
		"run": {
			Name:        "run",
			Description: "Run the 105chat server",
			FlagSet:     flag.NewFlagSet("run", flag.ExitOnError),
			Handler:     Run,
		},
		"migrate": {
			Name:        "migrate",
			Description: "Migrate database models",
			FlagSet:     flag.NewFlagSet("migrate", flag.ExitOnError),
			Handler:     Migrate,
		},
		"test": {
			Name:        "test",
			Description: "Run the unit tests for the API",
			FlagSet:     flag.NewFlagSet("test", flag.ExitOnError),
			Handler:     Test,
		},
	}

	if len(os.Args) < 2 {
		Help([]string{})
		os.Exit(0)
	}

	command, found := commands[os.Args[1]]
	if !found {
		fmt.Printf(
			"❓ Command %s is not valid. Run help command to get the full list of commands\n",
			os.Args[1],
		)
		os.Exit(1)
	}

	command.Handler(command.FlagSet.Args())
}

// Help is a CLI command that gives the command list
func Help(_ []string) {
	println("───── 105chat CLI help ─────\n")

	println("COMMANDS")
	for _, command := range commands {
		fmt.Println(command.String())
	}
}

// Run is a CLI command that starts the 105chat API
func Run(_ []string) {
	color.HiCyan("\n  _____   ________   ________   ________   ___  ___   ________   _________   \n / __  \\ |\\   __  \\ |\\   ____\\ |\\   ____\\ |\\  \\|\\  \\ |\\   __  \\ |\\___   ___\\ \n|\\/_|\\  \\\\ \\  \\|\\  \\\\ \\  \\___|_\\ \\  \\___| \\ \\  \\\\\\  \\\\ \\  \\|\\  \\\\|___ \\  \\_| \n\\|/ \\ \\  \\\\ \\  \\\\\\  \\\\ \\_____  \\\\ \\  \\     \\ \\   __  \\\\ \\   __  \\    \\ \\  \\  \n     \\ \\  \\\\ \\  \\\\\\  \\\\|____|\\  \\\\ \\  \\____ \\ \\  \\ \\  \\\\ \\  \\ \\  \\    \\ \\  \\ \n      \\ \\__\\\\ \\_______\\ ____\\_\\  \\\\ \\_______\\\\ \\__\\ \\__\\\\ \\__\\ \\__\\    \\ \\__\\\n       \\|__| \\|_______||\\_________\\\\|_______| \\|__|\\|__| \\|__|\\|__|     \\|__|\n                       \\|_________|                                          \n")
	db.OpenDatabase()

	go ws.HandlePipeline()
	http.Handle("/v1/ws", websocket.Handler(ws.Server))

	go httpServer.Server()

	addr := os.Getenv("WS_ADDRESS") + ":" + os.Getenv("WS_PORT")
	log.Println("▶ WS server listening on", color.CyanString("ws://"+addr))
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Panicf("‼ WS server fatal error: %s", err.Error())
	}
}

// Migrate is a CLI command that performs a migration on the database to register models
func Migrate(_ []string) {
	db.OpenDatabase()
	log.Println("✏ Database migration started")
	db.Database.AutoMigrate(
		&db.User{},
		&db.Message{},
		&db.Room{},
		&db.Group{},
		&db.RoomPermission{},
		&db.GroupInheritance{},
	)
	log.Println(color.GreenString("✅ Database migration complete"))
}

func Test(_ []string) {
	println("───── 105chat tests ─────\n")

	log.Println(color.CyanString("⏩ Step 1: database migrate"))
	Migrate([]string{})

	log.Println(color.CyanString("⏩ Step 2: HTTP API tests"))
	cmd := exec.Command("go", "test", "./tests")
	stdout, _ := cmd.Output()

	println(string(stdout))

	log.Println(color.HiGreenString("✅ Tests completed. Make sure to check the output above to search for any error."))
}
