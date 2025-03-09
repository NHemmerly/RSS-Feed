package state

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strings"
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
		pubDate, err := time.Parse("2006-01-02", item.PubDate)
		if err != nil {
			pubDate, err = time.Parse(time.RFC1123, item.PubDate)
			if err != nil {
				fmt.Printf("could not parse date: %s: %v\n", item.PubDate, err)
				pubDate = time.Now()
			}
		}
		_, err = s.Db.CreatePost(context.Background(), database.CreatePostParams{
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: true},
			PublishedAt: sql.NullTime{Time: pubDate, Valid: true},
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique constraint") {
				fmt.Printf("Skipping duplicate post: %s\n", item.Title)
				continue
			}
			fmt.Printf("could not create post: %s: %v\n", item.Title, err)
			continue
		}
	}
	return nil
}

func selectFeedNameUrl(name string, link string) (string, string) {
	var feedUrl string
	var feedName string
	if _, err := url.ParseRequestURI(name); err == nil {
		feedUrl = name
		feedName = link
	} else if _, err := url.ParseRequestURI(link); err == nil {
		feedUrl = link
		feedName = name
	}
	return feedName, feedUrl
}
