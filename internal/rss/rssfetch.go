package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

var feedClient = &http.Client{
	Timeout: 10 * time.Second,
}

func requestFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("User-Agent", "gator")

	resp, err := feedClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching feed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		io.Copy(io.Discard, io.LimitReader(resp.Body, 8*1024*1024))
		return nil, fmt.Errorf("error status code %d", resp.StatusCode)
	}

	var feed RSSFeed
	decoder := xml.NewDecoder(resp.Body)
	if err := decoder.Decode(&feed); err != nil {
		return nil, fmt.Errorf("error decoding xml: %w", err)
	}

	io.Copy(io.Discard, io.LimitReader(resp.Body, 8*1024*1024))

	return &feed, nil
}

func GetFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	feed, err := requestFeed(ctx, feedURL)
	if err != nil {
		return nil, fmt.Errorf("error occurred during http request: %w", err)
	}

	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)

	for i, item := range feed.Channel.Item {
		feed.Channel.Item[i].Description = html.UnescapeString(item.Description)
		feed.Channel.Item[i].Title = html.UnescapeString(item.Title)
	}

	return feed, nil
}
