package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/SebastianGeroli/gator-boot-dev/internal/config"
	"github.com/SebastianGeroli/gator-boot-dev/internal/database"
	"github.com/google/uuid"
)

type commands struct {
	handlers map[string]func(*state, command) error
}

type command struct {
	name string
	args []string
}

type state struct {
	db     *database.Queries
	config *config.Config
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

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("username is required\n")
	}
	username := cmd.args[0]
	user, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return err
	}
	s.config.SetUser(username)
	fmt.Printf("Logged in as: %v\n", user.Name)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("username is required\n")
	}
	name := cmd.args[0]
	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	}
	_, err := s.db.CreateUser(context.Background(), userParams)
	if err != nil {
		return err
	}
	s.config.SetUser(name)
	fmt.Printf("The user: %v was created\n", name)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteAllUsers(context.Background())
	return err
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range users {
		if user.Name == s.config.UserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {

			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", feed)

	return nil
}
