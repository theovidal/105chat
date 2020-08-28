package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("💾 No .env file at the root - Ignoring")
	}

	makeCommands()

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
