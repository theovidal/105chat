package main

import (
	"flag"
	"fmt"
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
var commands map[string]Command

func makeCommands() {
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
		"cache": {
			Name:        "cache",
			Description: "Populate the cache with common data",
			FlagSet:     flag.NewFlagSet("cache", flag.ExitOnError),
			Handler:     Cache,
		},
		"test": {
			Name:        "test",
			Description: "Run the unit tests for the API",
			FlagSet:     flag.NewFlagSet("test", flag.ExitOnError),
			Handler:     Test,
		},
	}
}
