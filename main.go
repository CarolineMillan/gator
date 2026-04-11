package main

import (
	"fmt"
	"gator/internal/cli"
	"gator/internal/config"
	"os"
)

func main() {
	// read the config file
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	// store it in the state struct
	s := cli.NewState(&cfg)

	cmds := cli.NewCommands()

	err = cmds.Register("login", cli.HandlerLogin)

	// get a handle on the args
	args := os.Args

	if len(args) < 2 {
		fmt.Println("Error: command name required. Usage: gator <command> <args>")
		os.Exit(1)
	}

	// create a command struct from the args
	c := cli.NewCommand(args[1], args[2:])

	// run the command
	err = cmds.Run(s, c)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
