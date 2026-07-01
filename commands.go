package main

import (
	"errors"
	"fmt"

	"github.com/Robot-tim1/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if handler, ok := c.handlers[cmd.name]; ok {
		return handler(s, cmd)
	}
	return errors.New("error command not found")
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("the login handler expects a single argument, the username.")
	}

	err := s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return fmt.Errorf("error setting username: %w", err)
	}
	fmt.Print("User has been set\n")
	return nil
}
