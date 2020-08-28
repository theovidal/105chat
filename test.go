package main

import (
	"log"
	"os/exec"

	"github.com/fatih/color"
)

func Test(_ []string) {
	println("───── 105chat tests ─────\n")

	log.Println(color.CyanString("⏩ Step 1: database migrate"))
	Migrate(nil)

	log.Println(color.CyanString("⏩ Step 2: HTTP API tests"))
	cmd := exec.Command("go", "test", "./tests")
	stdout, _ := cmd.Output()

	println(string(stdout))

	log.Println(color.HiGreenString("✅ Tests completed. Make sure to check the output above to search for any error."))
}
