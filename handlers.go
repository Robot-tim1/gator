package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Robot-tim1/gator/internal/database"
	"github.com/Robot-tim1/gator/internal/rss"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("the login handler expects a single argument, the username.")
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
		return errors.New("the register handler expects a single argument, a name.")
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

func handlerAgg(s *state, cmd command) error {
	feed, err := rss.GetFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("error occurred while getting feed: %w", err)
	}

	fmt.Println(feed)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return errors.New("The addfeed handler expects 2 arguments, the feed's name and url")
	}

	user, err := s.db.GetUserFromName(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error getting current user: %w", err)
	}

	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), params)
	if err != nil {
		return fmt.Errorf("error creating feed record: %w", err)
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
