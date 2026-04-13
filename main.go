package main

import _ "github.com/lib/pq"

import (
	"database/sql"
	"fmt"
	"gator/internal/cli"
	"gator/internal/config"
	"gator/internal/database"
	"os"
)

func main() {
	// read the config file
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	// open a connection to the database using the config url
	db, err := sql.Open("postgres", cfg.DBurl)

	// create a new *database.Queries
	dbQueries := database.New(db)

	// store it in the state struct
	s := cli.NewState(dbQueries, &cfg)

	cmds := cli.NewCommands()

	err = cmds.Register("login", cli.HandlerLogin)
	err = cmds.Register("register", cli.HandlerRegister)
	err = cmds.Register("reset", cli.HandlerReset)

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
