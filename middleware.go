package main

import (
	"context"
	"fmt"

	"github.com/Robot-tim1/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, c command) error {
		user, err := s.db.GetUserFromName(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("error getting current user: %w", err)
		}

		return handler(s, c, user)
	}
}
