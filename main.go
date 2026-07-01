package main

import (
	"fmt"
	"os"

	"github.com/Robot-tim1/gator/internal/config"
)

func main() {

	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("%v", err)
	}

	s := &state{cfg: &cfg}
	commands := commands{handlers: make(map[string]func(*state, command) error)}

	commands.register("login", handlerLogin)

	if len(os.Args) < 2 {
		fmt.Printf("Not enough arguments\n")
		os.Exit(1)
	}
	commandName := os.Args[1]
	args := os.Args[2:]
	cmd := command{name: commandName, args: args}
	err = commands.run(s, cmd)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
