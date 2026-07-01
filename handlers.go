package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Robot-tim1/gator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("the login handler expects a single argument, the username.")
	}

	name := cmd.args[0]

	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("error couldn't find user: %w", err)
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("error setting username: %w", err)
	}
	fmt.Println("User has been set")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("the register handler expects a single argument, a name.")
	}

	name := cmd.args[0]

	_, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name})
	if err != nil {
		return fmt.Errorf("error creating user record: %w", err)
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("error setting user: %w", err)
	}
	fmt.Println("User was created")
	return nil
}
