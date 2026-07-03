package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/Robot-tim1/gator/internal/config"
	"github.com/Robot-tim1/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}
	defer db.Close()

	s := &state{
		cfg: &cfg,
		db:  database.New(db),
	}

	commands := commands{handlers: make(map[string]func(*state, command) error)}

	commands.registerAll()

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	cmd := command{name: cmdName, args: cmdArgs}
	err = commands.run(s, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
