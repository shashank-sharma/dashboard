package models

import (
	"encoding/json"
	"fmt"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

const (
	SourceTypeRSS        = "rss"
	SourceTypeHackerNews = "hackernews"
	SourceTypeGithub     = "github"
	SourceTypeReddit     = "reddit"
	SourceTypeYoutube    = "youtube"
	SourceTypeCustom     = "custom"
)

// Feed categories
const (
	CategoryExplore   = "explore"
	CategoryFollowing = "following"
	CategoryToRead    = "toread"
)

// Status values
const (
	StatusUnread    = "unread"
	StatusRead      = "read"
	StatusSaved     = "saved"
	StatusDismissed = "dismissed"
)

var _ core.Model = (*FeedSource)(nil)
var _ core.Model = (*FeedItem)(nil)
var _ core.Model = (*FeedCategory)(nil)


type FeedSource struct {
	BaseModel

	User        string                  `db:"user" json:"user"`
	Name        string                  `db:"name" json:"name"`
	Type        string                  `db:"type" json:"type"`
	URL         string                  `db:"url" json:"url"`
	Config      string                  `db:"config" json:"config"`
	CategoryIDs types.JSONArray[string] `db:"category_ids" json:"category_ids"`
	RefreshRate int                     `db:"refresh_rate" json:"refresh_rate"`
	IsActive    bool                    `db:"is_active" json:"is_active"`
	LastFetched types.DateTime          `db:"last_fetched" json:"last_fetched"`
	ErrorCount  int                     `db:"error_count" json:"error_count"`
	LastError   string                  `db:"last_error" json:"last_error"`
}

// GetConfigMap converts the Config string to a map[string]interface{}
func (m *FeedSource) GetConfigMap() (map[string]interface{}, error) {
	if m.Config == "" {
		return map[string]interface{}{}, nil
	}

	var configMap map[string]interface{}
	if err := json.Unmarshal([]byte(m.Config), &configMap); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %w", err)
	}

	return configMap, nil
}

// SetConfigMap sets the Config string from a map[string]interface{}
func (m *FeedSource) SetConfigMap(configMap map[string]interface{}) error {
	if configMap == nil {
		m.Config = "{}"
		return nil
	}

	configBytes, err := json.Marshal(configMap)
	if err != nil {
		return fmt.Errorf("failed to marshal config to JSON: %w", err)
	}

	m.Config = string(configBytes)
	return nil
}


type FeedItem struct {
	BaseModel

	User        string                  `db:"user" json:"user"`
	SourceID    string                  `db:"source_id" json:"source_id"`
	ExternalID  string                  `db:"external_id" json:"external_id"`
	Title       string                  `db:"title" json:"title"`
	Content     string                  `db:"content" json:"content"`
	URL         string                  `db:"url" json:"url"`
	Author      string                  `db:"author" json:"author"`
	PublishedAt types.DateTime          `db:"published_at" json:"published_at"`
	FetchedAt   types.DateTime          `db:"fetched_at" json:"fetched_at"`
	Status      string                  `db:"status" json:"status"`             // unread, read, saved, dismissed
	Rating      int                     `db:"rating" json:"rating"`             // User rating (1-5)
	Tags        types.JSONArray[string] `db:"tags" json:"tags"`                 // Tags for the feed item (many-to-many)
	CategoryIDs types.JSONArray[string] `db:"category_ids" json:"category_ids"` // Categories this item belongs to
	Summary     string                  `db:"summary" json:"summary"`           // AI-generated summary
	Metadata    string                  `db:"metadata" json:"metadata"`         // Additional source-specific data
	IsProcessed bool                    `db:"is_processed" json:"is_processed"` // Whether AI processing is complete
}

// GetMetadataMap converts the Metadata string to a map[string]interface{}
func (m *FeedItem) GetMetadataMap() (map[string]interface{}, error) {
	if m.Metadata == "" {
		return map[string]interface{}{}, nil
	}

	var metadataMap map[string]interface{}
	if err := json.Unmarshal([]byte(m.Metadata), &metadataMap); err != nil {
		return nil, fmt.Errorf("failed to parse metadata JSON: %w", err)
	}

	return metadataMap, nil
}

// SetMetadataMap sets the Metadata string from a map[string]interface{}
func (m *FeedItem) SetMetadataMap(metadataMap map[string]interface{}) error {
	if metadataMap == nil {
		m.Metadata = "{}"
		return nil
	}

	metadataBytes, err := json.Marshal(metadataMap)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata to JSON: %w", err)
	}

	m.Metadata = string(metadataBytes)
	return nil
}

type FeedCategory struct {
	BaseModel

	User        string `db:"user" json:"user"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
	Color       string `db:"color" json:"color"`
	Icon        string `db:"icon" json:"icon"`
	Type        string `db:"type" json:"type"` // explore, following, toread, custom
	SortOrder   int    `db:"sort_order" json:"sort_order"`
	IsDefault   bool   `db:"is_default" json:"is_default"`
}

func (m *FeedSource) TableName() string {
	return "feed_sources"
}

func (m *FeedItem) TableName() string {
	return "feed_items"
}

func (m *FeedCategory) TableName() string {
	return "feed_categories"
}
