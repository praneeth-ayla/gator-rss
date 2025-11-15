package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/praneeth-ayla/go-rss/internal/database"
)

// RSSFeed represents the structure of an RSS feed XML.
type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

// RSSItem represents an individual item within an RSS feed.
type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

// fetchFeed fetches and parses an RSS feed from a given URL.
func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check for HTTP errors.
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	// Read the response body.
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var feed RSSFeed
	// Unmarshal XML data into the RSSFeed struct.
	err = xml.Unmarshal(data, &feed)
	if err != nil {
		return nil, err
	}

	// Unescape HTML entities in feed title and description.
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	// Unescape HTML entities in each item's title and description.
	for _, item := range feed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
	}

	return &feed, nil
}

// middlewareLoggedIn is a middleware that ensures a user is logged in before executing the handler.
func middlewareLoggedIn(
	handler func(s *state, cmd command, user database.User) error,
) func(*state, command) error {

	return func(s *state, cmd command) error {
		// Attempt to retrieve the current user from the database.
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return err // Return error if user is not found or other DB error.
		}

		// Execute the original handler with the retrieved user.
		return handler(s, cmd, user)
	}
}

// scrapeFeeds fetches the next feed to scrape, marks it as fetched, and processes its items.
func scrapeFeeds(s *state) error {
	ctx := context.Background()
	// Get the next feed that needs to be fetched.
	feedToFetch, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return err
	}

	// Mark the feed as fetched in the database.
	err = s.db.MarkFeedFetched(ctx, database.MarkFeedFetchedParams{
		ID:            feedToFetch.ID,
		UpdatedAt:     time.Now(),
		LastFetchedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})
	if err != nil {
		return err
	}

	// Fetch the content of the feed URL.
	feed, err := fetchFeed(ctx, feedToFetch.Url)
	if err != nil {
		return err
	}
	// Print titles and descriptions of feed items.
	for _, item := range feed.Channel.Item {
		fmt.Printf("Title: %v\n", item.Title)
		fmt.Printf("Description: %v\n", item.Description)
	}
	return nil
}
