package providers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/shashank-sharma/backend/internal/logger"
	"github.com/shashank-sharma/backend/internal/models"
	"github.com/shashank-sharma/backend/internal/services"
)

type RSSProvider struct {
	client *http.Client
	parser *gofeed.Parser
}

func NewRSSProvider() services.FeedSourceProvider {
	return &RSSProvider{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		parser: gofeed.NewParser(),
	}
}

func (p *RSSProvider) GetProviderType() string {
	return models.SourceTypeRSS
}

func (p *RSSProvider) FetchItems(ctx context.Context, config map[string]interface{}, lastFetched time.Time) ([]services.RawFeedItem, error) {
	url, ok := config["url"].(string)
	if !ok || url == "" {
		return nil, fmt.Errorf("RSS feed URL is required")
	}

	logger.LogInfo(fmt.Sprintf("Fetching RSS feed: %s", url))

	// Parse the feed
	feed, err := p.parser.ParseURLWithContext(url, ctx)
	if err != nil {
		return nil, fmt.Errorf("error parsing RSS feed: %v", err)
	}

	// Process feed items
	var items []services.RawFeedItem
	for _, item := range feed.Items {
		// Skip items older than the last fetch time
		if !lastFetched.IsZero() && !item.PublishedParsed.IsZero() && !item.PublishedParsed.After(lastFetched) {
			continue
		}

		// Create feed item
		feedItem := services.RawFeedItem{
			ExternalID: item.GUID,
			Title:      item.Title,
			Content:    item.Description,
			URL:        item.Link,
			Metadata:   make(map[string]interface{}),
		}

		if item.Author != nil {
			feedItem.Author = item.Author.Name
		}

		if item.PublishedParsed != nil {
			feedItem.PublishedAt = *item.PublishedParsed
		} else {
			feedItem.PublishedAt = time.Now()
		}

		if len(item.Categories) > 0 {
			feedItem.Tags = item.Categories
		}

		if feed.Title != "" {
			feedItem.Metadata["feed_title"] = feed.Title
		}
		if feed.Link != "" {
			feedItem.Metadata["feed_link"] = feed.Link
		}
		if item.Content != "" {
			feedItem.Metadata["full_content"] = item.Content
		}

		items = append(items, feedItem)
	}

	logger.LogInfo(fmt.Sprintf("Fetched %d items from RSS feed", len(items)))
	return items, nil
}

func (p *RSSProvider) Validate(config map[string]interface{}) error {
	url, ok := config["url"].(string)
	if !ok || url == "" {
		return fmt.Errorf("RSS feed URL is required")
	}

	// Try fetching the feed to validate it
	_, err := p.parser.ParseURL(url)
	if err != nil {
		return fmt.Errorf("invalid RSS feed URL: %v", err)
	}

	return nil
}
