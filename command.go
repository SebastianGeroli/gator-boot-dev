package main

import (
	"errors"
	"fmt"

	"github.com/SebastianGeroli/gator-boot-dev/internal/config"
)

type commands struct {
	handlers map[string]func(*state, command) error
}

type command struct {
	name string
	args []string
}

type state struct {
	config *config.Config
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("username is required\n")
	}
	username := cmd.args[0]
	err := s.config.SetUser(username)
	if err != nil {
		return err
	}
	fmt.Printf("The user has been set\n")
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	cmdValue, ok := c.handlers[cmd.name]
	if !ok {
		return errors.New("Command does not exist\n")
	}
	err := cmdValue(s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}
