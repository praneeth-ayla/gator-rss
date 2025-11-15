package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/praneeth-ayla/gator/internal/config"
	"github.com/praneeth-ayla/gator/internal/database"
)

// state holds the application's global state, including database queries and configuration.
type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	// Read application configuration.
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	// Initialize program state.
	programState := &state{
		cfg: &cfg,
	}

	// Open database connection.
	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		log.Fatal("Error connecting to db:", err)
	}

	// Create database queries instance.
	dbQueries := database.New(db)
	programState.db = dbQueries

	// Initialize commands and register handlers.
	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerGetUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", (handlerFeeds))
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))

	// Check for command-line arguments.
	if len(os.Args) < 2 {
		log.Fatal("Usage: cli  [args...]")
	}
	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	// Run the specified command.
	err = cmds.run(programState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		os.Exit(1)
		log.Fatal(err)
	}

}
