package main

import (
	"context"

	"github.com/SebastianGeroli/gator-boot-dev/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUserByName(context.Background(), s.config.UserName)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
}
