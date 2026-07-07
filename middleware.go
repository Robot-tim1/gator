package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Robot-tim1/gator/internal/database"
	"github.com/Robot-tim1/gator/internal/rss"
	"github.com/google/uuid"
	"github.com/lib/pq"
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

const Layout = "Mon, 02 Jan 2006 15:04:05 -0700"

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

		pubDate, err := time.Parse(Layout, item.PubDate)
		if err != nil {
			return fmt.Errorf("error parsing publication date: %w", err)
		}

		description := sql.NullString{String: "", Valid: false}
		if item.Description != "" {
			description = sql.NullString{String: item.Description, Valid: true}
		}

		CreatePostParams := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: description,
			PublishedAt: pubDate,
			FeedID:      nextFeed.ID,
		}

		_, err = s.db.CreatePost(context.Background(), CreatePostParams)
		if err != nil {
			var pqErr *pq.Error
			if errors.As(err, &pqErr) && pqErr.Code == "23505" {
				// do nothing
			} else {
				return fmt.Errorf("error creating post record: %w", err)
			}
		}
	}
	return nil
}
