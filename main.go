package main

import (
	"fmt"
	"os"

	"github.com/SebastianGeroli/gator-boot-dev/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
	appState := state{
		config: &cfg,
	}
	appCommands := commands{
		handlers: map[string]func(*state, command) error{
			"login": handlerLogin,
		},
	}

	args := os.Args
	if len(args) < 2 {
		fmt.Printf("insufficient arguments\n")
		os.Exit(1)
	}

	cmdName := args[1]
	cmdArgs := args[2:]
	cmd := command{
		name: cmdName,
		args: cmdArgs,
	}
	err = appCommands.run(&appState, cmd)
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
}
