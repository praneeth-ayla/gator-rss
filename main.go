package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/praneeth-ayla/go-rss/internal/config"
)

type state struct {
	cfg *config.Config
}

func handlerLogin(s *state, cmd command) error {

	if len(cmd.Args) < 1 {
		return errors.New("the login handler expects a single argument, the username")
	}

	username := cmd.Args[0]
	err := s.cfg.SetUser(username)
	if err != nil {
		return err
	}

	fmt.Printf("the user %s has been set", username)
	return nil
}

func main() {
	cfg, err := config.Read()

	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	programState := &state{
		cfg: &cfg,
	}

	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}
	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	err = cmds.run(programState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}
}
