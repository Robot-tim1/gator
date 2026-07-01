package main

import (
	"errors"
)

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
