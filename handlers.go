package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Robot-tim1/gator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("the login command expects a single argument, a name.")
	}

	name := cmd.args[0]

	_, err := s.db.GetUserFromName(context.Background(), name)
	if err != nil {
		return fmt.Errorf("error couldn't find user: %w", err)
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("error setting username: %w", err)
	}
	fmt.Println("User has been set")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("the register command expects a single argument, a name.")
	}

	name := cmd.args[0]
	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	}

	_, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		return fmt.Errorf("error creating user record: %w", err)
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("error setting user: %w", err)
	}
	fmt.Println("User was created")
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error deleting users: %w", err)
	}
	fmt.Println("Users table has been reset")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error getting users: %w", err)
	}

	for _, user := range users {
		if s.cfg.CurrentUserName == user {
			fmt.Printf("* %s (current)\n", user)
		} else {
			fmt.Printf("* %s\n", user)
		}
	}
	return nil
}

func handlerAgg(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return errors.New("The agg command expects on argument, the time between requests like 1s, 1m, 1h, etc.")
	}

	time_between_reqs := cmd.args[0]

	timeReq := time.Second * 10
	inputTime, err := time.ParseDuration(time_between_reqs)
	if err != nil {
		fmt.Printf("error parsing input time, using default time of 10 seconds\n")
	} else {
		timeReq = inputTime
	}

	fmt.Printf("Collecting feeds every %s\n", timeReq)
	ticker := time.Tick(timeReq)
	err = scrapeFeeds(s, user)
	if err != nil {
		log.Printf("error scraping feeds: %v", err)
	}
	for range ticker {
		err = scrapeFeeds(s, user)
		if err != nil {
			log.Printf("error scraping feeds: %v", err)
		}
	}
	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return errors.New("The addfeed command expects 2 arguments, the feed's name and url")
	}

	createParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), createParams)
	if err != nil {
		return fmt.Errorf("error creating feed record: %w", err)
	}

	followParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	_, err = s.db.CreateFeedFollow(context.Background(), followParams)
	if err != nil {
		return fmt.Errorf("error following created feed: %w", err)
	}

	fmt.Println(feed)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error getting feeds from database: %w", err)
	}

	for _, feed := range feeds {
		user, err := s.db.GetUserFromId(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("error getting user of feed: %w", err)
		}

		fmt.Printf("(%s: %s) feed created by %s\n", feed.Name, feed.Url, user.Name)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return errors.New("The follow command expects on argument, the url of a feed")
	}

	feed, err := s.db.GetFeedFromUrl(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("error getting feed from database: %w", err)
	}

	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	followRow, err := s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return fmt.Errorf("error following feed: %w", err)
	}

	fmt.Printf("%s followed %s\n", followRow.UserName, followRow.FeedName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	followedFeeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error getting followed feeds for user: %w", err)
	}

	if len(followedFeeds) == 0 {
		return errors.New("You're not following any feeds")
	}

	fmt.Printf("%s is following\n", followedFeeds[0].UserName)
	for _, item := range followedFeeds {
		fmt.Printf(" - %s\n", item.FeedName)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return errors.New("The unfollow command expects on argument, the url of a feed")
	}

	feed, err := s.db.GetFeedFromUrl(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("error getting feed from database: %w", err)
	}

	params := database.DeleteFeedFollowParams{UserID: user.ID, FeedID: feed.ID}

	err = s.db.DeleteFeedFollow(context.Background(), params)
	if err != nil {
		return fmt.Errorf("error deleting from feed_follow table: %w", err)
	}

	fmt.Printf("%s unfollowed %s\n", user.Name, feed.Name)
	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int32
	limit = 2
	if len(cmd.args) == 1 {
		templimit, err := strconv.ParseInt(cmd.args[0], 10, 64)
		if err == nil {
			limit = int32(templimit)
		}
	}

	getPostsParams := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	}

	posts, err := s.db.GetPostsForUser(context.Background(), getPostsParams)
	if err != nil {
		return fmt.Errorf("error getting followed posts: %w", err)
	}

	for _, post := range posts {
		fmt.Printf("Title: %s\n", post.Title)
		fmt.Printf("Published at %s\n\n", post.PublishedAt)
		fmt.Printf("Description: %s\n\n", post.Description.String)
		fmt.Printf("Link: %s\n\n", post.Url)
	}
	return nil
}
