package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/NHemmerly/RSS-Feed/internal/config"
	"github.com/NHemmerly/RSS-Feed/internal/database"
	"github.com/NHemmerly/RSS-Feed/internal/state"
	_ "github.com/lib/pq"
)

func main() {
	userConfig, err := config.Read()
	if err != nil {
		fmt.Printf("error reading config\n")
	}
	var userState state.State
	userState.Cfg = userConfig
	var cmds state.Commands
	cmds.Cmds = map[string]func(*state.State, state.Command) error{}
	db, err := sql.Open("postgres", userConfig.DbURL)
	if err != nil {
		fmt.Printf("could not open database connection")
		os.Exit(1)
	}
	dbQueries := database.New(db)

	userState.Db = dbQueries

	args := os.Args
	cmds.Register("login", state.HandlerLogin)
	cmds.Register("register", state.HandlerRegister)
	cmds.Register("reset", state.HandlerReset)
	cmds.Register("users", state.HandlerUsers)
	cmds.Register("agg", state.HandlerAgg)
	cmds.Register("addfeed", state.HandlerAddFeed)
	cmds.Register("feeds", state.HandlerFeeds)
	cmds.Register("follow", state.HandlerFollow)
	cmds.Register("following", state.HandlerFollowing)
	if len(args) < 2 {
		fmt.Println("Why two?")
		os.Exit(1)
	}
	command := state.Command{Name: args[1], Args: args[2:]}
	if err = cmds.Run(&userState, command); err != nil {
		fmt.Printf("could not run command: %s", err)
		os.Exit(1)
	}

}
