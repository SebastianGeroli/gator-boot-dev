package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/SebastianGeroli/gator-boot-dev/internal/config"
	"github.com/SebastianGeroli/gator-boot-dev/internal/database"
	_ "github.com/lib/pq"
)

func main() {

	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	dbQueries := database.New(db)

	appState := state{
		db:     dbQueries,
		config: &cfg,
	}
	appCommands := commands{handlers: map[string]func(*state, command) error{}}
	appCommands.register("login", handlerLogin)
	appCommands.register("register", handlerRegister)
	appCommands.register("reset", handlerReset)
	appCommands.register("users", handlerUsers)
	appCommands.register("agg", handlerAgg)
	appCommands.register("addfeed", handlerAddFeed)
	appCommands.register("feeds", handlerFeeds)
	appCommands.register("follow", handlerFollow)
	appCommands.register("following", handlerFollowing)

	args := os.Args
	if len(args) < 2 {
		fmt.Printf("insufficient arguments\n")
		os.Exit(1)
	}

	cmdName := args[1]
	cmdArgs := args[2:]
	cmd := command{
		name: cmdName,
		args: cmdArgs,
	}
	err = appCommands.run(&appState, cmd)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
