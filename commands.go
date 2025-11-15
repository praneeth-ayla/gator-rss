package main

import "errors"

// command represents a single command with its name and arguments.
type command struct {
	Name string
	Args []string
}

// commands manages the registration and execution of application commands.
type commands struct {
	registeredCommands map[string]func(*state, command) error
}

// register adds a command handler to the map.
func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}

// run executes a command by its name with the given state and arguments.
func (c *commands) run(s *state, cmd command) error {
	f, ok := c.registeredCommands[cmd.Name]
	if !ok {
		return errors.New("command not found")
	}
	return f(s, cmd)
}
