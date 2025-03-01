package services

import (
	"context"
	"time"

	"github.com/shashank-sharma/backend/internal/models"
)

// FeedItem represents a raw item fetched from a source
// before being processed and stored
type RawFeedItem struct {
	ExternalID  string
	Title       string
	Content     string
	URL         string
	Author      string
	PublishedAt time.Time
	Tags        []string
	Metadata    map[string]interface{}
}

// FeedSourceProvider defines the interface for all feed source implementations
type FeedSourceProvider interface {
	// GetProviderType returns the type of feed source (rss, hackernews, etc.)
	GetProviderType() string

	// FetchItems fetches items from the source
	// config contains source-specific configuration
	FetchItems(ctx context.Context, config map[string]interface{}, lastFetched time.Time) ([]RawFeedItem, error)

	// Validate validates the source configuration
	Validate(config map[string]interface{}) error
}

// FeedProcessor handles processing of raw feed items
type FeedProcessor interface {
	// ProcessItem processes a raw feed item and returns a processed feed item
	ProcessItem(ctx context.Context, item *models.FeedItem) error

	// GenerateSummary generates an AI summary for a feed item
	GenerateSummary(ctx context.Context, item *models.FeedItem) (string, error)

	// SuggestTags suggests tags for a feed item based on content
	SuggestTags(ctx context.Context, item *models.FeedItem) ([]string, error)
}

// FeedService coordinates the feed system
type FeedService interface {
	// RegisterProvider registers a new feed source provider
	RegisterProvider(provider FeedSourceProvider)

	// GetProvider returns the provider for a given source type
	GetProvider(sourceType string) (FeedSourceProvider, bool)

	// FetchFromSource fetches items from a specific source
	FetchFromSource(ctx context.Context, source *models.FeedSource) error

	// FetchAllSources fetches items from all active sources
	FetchAllSources(ctx context.Context) error

	// ProcessNewItems processes all unprocessed items
	ProcessNewItems(ctx context.Context) error

	// GetUserFeeds gets feed items for a user with filtering options
	GetUserFeeds(ctx context.Context, userID string, options map[string]interface{}) ([]*models.FeedItem, error)
}
