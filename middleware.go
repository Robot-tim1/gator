package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Robot-tim1/gator/internal/database"
	"github.com/Robot-tim1/gator/internal/rss"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, c command) error {
		user, err := s.db.GetUserFromName(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("error getting current user: %w", err)
		}

		return handler(s, c, user)
	}
}

// I know this function doesn't really fit here
// but I've got no other place to put it
func scrapeFeeds(s *state, user database.User) error {
	userFeedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error getting users follows: %w", err)
	}

	if len(userFeedFollows) == 0 {
		return errors.New("user is not following any feeds to aggregate")
	}

	nextFeed, err := s.db.GetNextFeedFromFollows(context.Background(), userFeedFollows)
	if err != nil {
		return fmt.Errorf("error getting next feed: %w", err)
	}

	markParams := database.MarkFeedFetchedParams{
		UpdatedAt: time.Now(),
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		ID: nextFeed.ID,
	}
	s.db.MarkFeedFetched(context.Background(), markParams)

	fetchedFeed, err := rss.GetFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return fmt.Errorf("error getting feed: %w", err)
	}

	for _, item := range fetchedFeed.Channel.Item {
		fmt.Printf("%s\n", item.Title)
	}
	return nil
}
