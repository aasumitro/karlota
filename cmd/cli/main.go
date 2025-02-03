package main

import (
	"log"

	"karlota.aasumitro.id/cmd/cli/command"
)

func main() {
	if err := command.CliCommand.Execute(); err != nil {
		log.Fatalf("Error executing CLI command: %s", err.Error())
	}
}
