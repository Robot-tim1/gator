package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
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

const Layout1 = "Mon, 02 Jan 2006 15:04:05 -0700"
const Layout2 = "Mon, 02 Jan 2006 15:04:05 MST"

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

	for _, followedFeed := range userFeedFollows {
		go scraper(s, followedFeed)
	}
	return nil
}

func scraper(s *state, followedFeed database.GetFeedFollowsForUserRow) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in scraper: %v", r)
		}
	}()

	markParams := database.MarkFeedFetchedParams{
		UpdatedAt: time.Now(),
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		ID: followedFeed.FeedID,
	}
	err := s.db.MarkFeedFetched(context.Background(), markParams)
	if err != nil {
		log.Printf("Couldn't mark feed %s as fetched: %v", followedFeed.FeedID, err)
		return
	}

	fetchedFeed, err := rss.GetFeed(context.Background(), followedFeed.FeedUrl)
	if err != nil {
		log.Printf("Error getting feed %s: %v", followedFeed.FeedUrl, err)
		return
	}
	log.Printf("fetched from feed: %s", fetchedFeed.Channel.Title)
	layouts := []string{Layout1, Layout2, time.ANSIC, time.RFC3339, time.RFC822}
	for _, item := range fetchedFeed.Channel.Item {
		var pubDate time.Time
		for _, layout := range layouts {
			pubDate, err = time.Parse(layout, item.PubDate)
			if err == nil {
				break
			}
		}
		if pubDate.IsZero() {
			log.Printf("Error parsing date for item %s: no matching layouts", item.Title)
			continue
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
			FeedID:      followedFeed.FeedID,
		}

		_, err = s.db.CreatePost(context.Background(), CreatePostParams)
		if err != nil {
			var pqErr *pq.Error
			if errors.As(err, &pqErr) && pqErr.Code == "23505" {
				continue
			}
			log.Printf("Error creating post record: %v", err)
		}
	}
}
