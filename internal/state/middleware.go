package state

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/NHemmerly/RSS-Feed/internal/database"
	"github.com/NHemmerly/RSS-Feed/internal/fetch"
)

func MiddlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		user, err := s.Db.GetUser(context.Background(), s.Cfg.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}

}

func scrapeFeeds(s *State) error {
	feed, err := s.Db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("could not fetch next feed: %w", err)
	}
	err = s.Db.MarkFeedFetched(context.Background(),
		database.MarkFeedFetchedParams{
			UpdatedAt:     time.Now(),
			LastFetchedAt: sql.NullTime{Time: time.Now(), Valid: true},
			ID:            feed.ID,
		})
	if err != nil {
		return fmt.Errorf("could not mark feed: %w", err)
	}
	feedData, err := fetch.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("could not fetch feed data: %w", err)
	}
	for _, item := range feedData.Channel.Item {
		fmt.Printf("%v - %v\n", item.Title, item.Link)
	}
	return nil
}
