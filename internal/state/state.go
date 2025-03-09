package state

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/NHemmerly/RSS-Feed/internal/config"
	"github.com/NHemmerly/RSS-Feed/internal/database"
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
	if len(cmd.Args) < 1 {
		return fmt.Errorf("not enough args, expected 1; ")
	}
	timeBetweenReqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("could not parse time argument: %w", err)
	}
	fmt.Printf("Collecting feeds every %v\n", timeBetweenReqs)
	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func HandlerAddFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("not enough args provided, expected 2; ")
	}
	feedName, feedUrl := selectFeedNameUrl(cmd.Args[0], cmd.Args[1])
	feed, err := s.Db.CreateFeed(context.Background(),
		database.CreateFeedParams{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      feedName,
			Url:       feedUrl,
			UserID:    user.ID,
		})
	if err != nil {
		return fmt.Errorf("could not create feed: %w", err)
	}
	_, err = s.Db.CreateFeedFollow(context.Background(),
		database.CreateFeedFollowParams{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			FeedID:    feed.ID,
		})
	if err != nil {
		return fmt.Errorf("could not create feed follow: %w", err)
	}
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

func HandlerFollow(s *State, cmd Command, user database.User) error {
	feedId, err := s.Db.GetFeedByUrl(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("could not get feed by url: %w", err)
	}
	if len(cmd.Args) < 1 {
		return fmt.Errorf("not enough args provided, expected 1; ")
	}
	s.Db.CreateFeedFollow(context.Background(),
		database.CreateFeedFollowParams{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			FeedID:    feedId,
		})
	return nil
}

func HandlerFollowing(s *State, cmd Command, user database.User) error {
	userFeeds, err := s.Db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		return fmt.Errorf("could not retrieve user's feeds: %w", err)
	}
	for i := range userFeeds {
		fmt.Println(userFeeds[i])
	}
	return nil
}

func HandlerUnfollow(s *State, cmd Command, user database.User) error {
	err := s.Db.RemoveFeedFollow(context.Background(),
		database.RemoveFeedFollowParams{
			Name: user.Name,
			Url:  cmd.Args[0],
		})
	if err != nil {
		return fmt.Errorf("could not remove follow record: %w ", err)
	}
	return nil
}

func userExists(s *State, name string) bool {
	_, err := s.Db.GetUser(context.Background(), name)
	return err == nil
}
