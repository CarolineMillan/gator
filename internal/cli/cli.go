package cli

import (
	"errors"
	"fmt"
	"gator/internal/config"
)

type state struct {
	cfg *config.Config
}

func NewState(cfg *config.Config) *state {
	s := state{}
	s.cfg = cfg
	return &s
}

type command struct {
	name string
	args []string
}

func NewCommand(name string, args []string) command {
	c := command{}
	c.name = name
	c.args = args
	return c
}

type commands struct {
	Handlers map[string]func(*state, command) error // a map of command names to their handler functions
}

func NewCommands() *commands {
	c := commands{}
	c.Handlers = make(map[string]func(*state, command) error)
	return &c
}

func (c *commands) Run(s *state, cmd command) error {
	// runs a given command with the provided state if it exists
	f, exists := c.Handlers[cmd.name]
	if !exists {
		return fmt.Errorf("Error: command %s does not exist.", cmd.name)
	}
	err := f(s, cmd)
	return err
}

func (c *commands) Register(name string, f func(*state, command) error) error {
	// registers a new handler function for a command name, i.e. adds it to the commands struct
	c.Handlers[name] = f
	return nil
}

func HandlerLogin(s *state, c command) error {
	if len(c.args) != 1 {
		return errors.New("login takes one argument. Usage: gator login <username>")
	}

	err := s.cfg.SetUser(c.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("User %s has been set successfully.", c.args[0])
	return nil
}
