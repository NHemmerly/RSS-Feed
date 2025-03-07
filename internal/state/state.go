package state

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/NHemmerly/RSS-Feed/internal/config"
	"github.com/NHemmerly/RSS-Feed/internal/database"
	"github.com/NHemmerly/RSS-Feed/internal/fetch"
	"github.com/google/uuid"
)

type State struct {
	Db  *database.Queries
	Cfg *config.Config
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Cmds map[string]func(*State, Command) error
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.Cmds[name] = f
}

func (c *Commands) Run(s *State, cmd Command) error {
	if val, ok := c.Cmds[cmd.Name]; ok {
		return val(s, cmd)
	} else {
		return fmt.Errorf("could not find function in cmds")
	}
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("no args provided; expected 1")
	}
	if !userExists(s, cmd.Args[0]) {
		return fmt.Errorf("user does not exist; ")
	}
	s.Cfg.SetUser(cmd.Args[0])
	fmt.Printf("User %s has been set!\n", cmd.Args[0])
	return nil

}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("no args provided; expected 1")
	}
	if !userExists(s, cmd.Args[0]) {
		user, err := s.Db.CreateUser(context.Background(),
			database.CreateUserParams{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Name:      cmd.Args[0]},
		)
		if err != nil {
			return fmt.Errorf("could not create user: %w", err)
		}
		fmt.Printf("User %s created! %v\n", user.Name, user)
	} else {
		fmt.Printf("User %s already exists\n", cmd.Args[0])
		os.Exit(1)
	}

	s.Cfg.SetUser(cmd.Args[0])
	return nil
}

func HandlerReset(s *State, cmd Command) error {
	if err := s.Db.ResetUser(context.Background()); err != nil {
		return fmt.Errorf("could not reset table: %w", err)
	}
	if err := s.Db.ResetFeed(context.Background()); err != nil {
		return fmt.Errorf("could not reset table: %w", err)
	}
	return nil
}

func HandlerUsers(s *State, cmd Command) error {
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("could not retrieve all users: %w", err)
	}
	for _, user := range users {
		if user.Name == s.Cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
			continue
		}
		fmt.Printf("* %s\n", user.Name)
	}

	return nil
}

func HandlerAgg(s *State, cmd Command) error {
	feed, err := fetch.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("could not fetch feed: %w", err)
	}
	fmt.Printf("%v", feed)
	return nil
}

func HandlerAddFeed(s *State, cmd Command) error {
	var feedUrl string
	var feedName string
	currentUser, err := s.Db.GetUser(context.Background(), s.Cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("could not find current user: %w", err)
	}
	if len(cmd.Args) < 2 {
		return fmt.Errorf("not enough args provided, expected 2; ")
	}
	if _, err := url.ParseRequestURI(cmd.Args[0]); err == nil {
		feedUrl = cmd.Args[0]
		feedName = cmd.Args[1]
	} else if _, err := url.ParseRequestURI(cmd.Args[1]); err == nil {
		feedUrl = cmd.Args[1]
		feedName = cmd.Args[0]
	}
	s.Db.CreateFeed(context.Background(),
		database.CreateFeedParams{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      feedName,
			Url:       feedUrl,
			UserID:    currentUser.ID,
		})
	return nil
}

func HandlerFeeds(s *State, cmd Command) error {
	feeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("could not get feeds: %w", err)
	}
	for i := range feeds {
		fmt.Println(feeds[i])
	}
	return nil
}

func userExists(s *State, name string) bool {
	_, err := s.Db.GetUser(context.Background(), name)
	return err == nil
}
