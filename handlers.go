package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/praneeth-ayla/go-rss/internal/database"
)

// handlerLogin handles user login by setting the current user in the config.
func handlerLogin(s *state, cmd command) error {

	if len(cmd.Args) < 1 {
		return errors.New("the login handler expects a single argument, the username")
	}

	username := cmd.Args[0]
	// Retrieve user from the database.
	user, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return err
	}

	// Set the current user in the application configuration.
	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Printf("the user %s has been set\n", username)
	return nil
}

// handlerRegister handles new user registration.
func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return errors.New("the register handler expects a single argument, the name")
	}
	name := cmd.Args[0]

	// Create a new user in the database.
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return err
	}

	// Set the newly created user as the current user.
	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}
	fmt.Println("the user was created", user)
	return nil
}

// handlerReset deletes all users from the database.
func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return err
	}
	return nil
}

// handlerGetUsers retrieves and prints all users, marking the current user.
func handlerGetUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	currentUser := s.cfg.CurrentUserName

	// Iterate and print each user, indicating the current one.
	for _, user := range users {
		if user.Name == currentUser {
			fmt.Printf("* %v (current)\n", user.Name)
			continue
		}
		fmt.Printf("* %v\n", user.Name)
	}

	return nil
}

// handlerAgg continuously scrapes feeds at a specified interval.
func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s ", cmd.Name)
	}

	// Parse the duration for time between requests.
	timeBetweenReqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}

	// Create a new ticker that fires at the specified interval.
	ticker := time.NewTicker(timeBetweenReqs)
	// Continuously scrape feeds on each tick.
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

// handlerAddFeed adds a new feed and creates a follow for the given user.
func handlerAddFeed(s *state, cmd command, user database.User) error {
	ctx := context.Background()

	name := cmd.Args[0]
	url := cmd.Args[1]
	if name == "" || url == "" {
		return errors.New("add feed command requires name and url")
	}

	// Create the feed in the database.
	feed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return err
	}

	// Create a feed follow for the user for this new feed.
	_, err = s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    feed.UserID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return err
	}

	fmt.Println(feed)

	return nil
}

// handlerFeeds retrieves and prints all registered feeds.
func handlerFeeds(s *state, cmd command) error {
	ctx := context.Background()
	feeds, err := s.db.GetFeeds(ctx)
	if err != nil {
		return err
	}

	// Iterate and print details for each feed.
	for _, feed := range feeds {
		user, err := s.db.GetUserById(ctx, feed.UserID)
		if err != nil {
			return err
		}
		fmt.Printf(`Feed Name: %v,
Feed URL: %v,
User Name: %v
`, feed.Name, feed.Url, user.Name)
	}

	return nil
}

// handlerFollow creates a feed follow for a given feed URL and current user.
func handlerFollow(s *state, cmd command, user database.User) error {
	ctx := context.Background()

	url := cmd.Args[0]
	// Get the feed by its URL.
	feed, err := s.db.GetFeedByUrl(ctx, url)
	if err != nil {
		return err
	}
	// Create the feed follow record.
	follow, err := s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return err
	}
	fmt.Printf(`Feed Name: %v
User Name: %v
`, follow.FeedName, follow.UserName)

	return nil
}

// handlerFollowing retrieves and prints all feeds followed by the current user.
func handlerFollowing(s *state, cmd command, user database.User) error {
	ctx := context.Background()

	feedFollows, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return nil
	}

	// Print the name of each followed feed.
	for _, feedFollow := range feedFollows {
		fmt.Println(feedFollow.FeedName)
	}

	return nil
}

// handlerUnfollow deletes a feed follow for a given URL and current user.
func handlerUnfollow(s *state, cmd command, user database.User) error {
	ctx := context.Background()
	url := cmd.Args[0]
	username := s.cfg.CurrentUserName

	// Delete the feed follow record.
	err := s.db.DeleteFeedFollow(ctx, database.DeleteFeedFollowParams{
		Url:  url,
		Name: username,
	})
	if err != nil {
		return err
	}

	return nil
}
