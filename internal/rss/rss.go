package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
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

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	/// fill out the RSSFeed struct, ie read in the RSS feed from the given url

	// http request to get the feed

	// this creates a request (doesn't send it)
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)

	if err != nil {
		return nil, fmt.Errorf("Error creating http request: %w", err)
	}

	// customise the request
	req.Header.Set("User-Agent", "gator")

	// send the request
	client := &http.Client{}
	res, err := client.Do(req)

	//res, err := http.Get(feedURL)
	if err != nil {
		return nil, fmt.Errorf("Error getting feed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("bad status: %s", res.Status)
	}

	// read into a []byte
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// unmarshall the xml data
	var feed RSSFeed
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		return nil, err
	}

	// unescape the html strings
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}

	return &feed, nil
}
