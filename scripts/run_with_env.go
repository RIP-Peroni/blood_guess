package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		log.Fatal("The .env file was not found. Create one based on .env.example")
	}

	cmd := exec.Command("go", "run", "cmd/bot/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
}
