package feed

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/shashank-sharma/backend/internal/logger"
	"github.com/shashank-sharma/backend/internal/models"
	"github.com/shashank-sharma/backend/internal/query"
	"github.com/shashank-sharma/backend/internal/services"
	"github.com/shashank-sharma/backend/internal/util"
)

type FeedServiceImpl struct {
	providers   map[string]services.FeedSourceProvider
	processor   services.FeedProcessor
	providersMu sync.RWMutex
}

func NewFeedService(processor services.FeedProcessor) services.FeedService {
	return &FeedServiceImpl{
		providers: make(map[string]services.FeedSourceProvider),
		processor: processor,
	}
}

// RegisterProvider registers a new feed source provider
func (s *FeedServiceImpl) RegisterProvider(provider services.FeedSourceProvider) {
	s.providersMu.Lock()
	defer s.providersMu.Unlock()

	providerType := provider.GetProviderType()
	s.providers[providerType] = provider
	logger.LogInfo(fmt.Sprintf("Registered feed provider: %s", providerType))
}

// GetProvider returns the provider for a given source type
func (s *FeedServiceImpl) GetProvider(sourceType string) (services.FeedSourceProvider, bool) {
	s.providersMu.RLock()
	defer s.providersMu.RUnlock()

	provider, exists := s.providers[sourceType]
	return provider, exists
}

// FetchFromSource fetches items from a specific source
func (s *FeedServiceImpl) FetchFromSource(ctx context.Context, source *models.FeedSource) error {
	provider, exists := s.GetProvider(source.Type)
	if !exists {
		return fmt.Errorf("no provider registered for source type: %s", source.Type)
	}

	logger.LogInfo(fmt.Sprintf("Fetching items from source: %s (ID: %s)", source.Name, source.Id))

	configMap, err := source.GetConfigMap()
	if err != nil {
		query.UpdateRecord[*models.FeedSource](source.Id, map[string]interface{}{
			"error_count":  source.ErrorCount + 1,
			"last_error":   fmt.Sprintf("Error parsing config: %v", err),
			"last_fetched": types.NowDateTime(),
		})
		return fmt.Errorf("error parsing config for source %s: %v", source.Name, err)
	}

	// Fetch items from the provider
	rawItems, err := provider.FetchItems(ctx, configMap, source.LastFetched.Time())
	if err != nil {
		query.UpdateRecord[*models.FeedSource](source.Id, map[string]interface{}{
			"error_count":  source.ErrorCount + 1,
			"last_error":   err.Error(),
			"last_fetched": types.NowDateTime(),
		})
		return fmt.Errorf("error fetching from source %s: %v", source.Name, err)
	}

	// Process and store each item
	for _, rawItem := range rawItems {
		existingItem, err := query.FindByFilter[*models.FeedItem](map[string]interface{}{
			"source_id":   source.Id,
			"external_id": rawItem.ExternalID,
		})

		// If the item exists, skip it
		if err == nil && existingItem != nil {
			continue
		}

		feedItem := &models.FeedItem{
			User:        source.User,
			SourceID:    source.Id,
			ExternalID:  rawItem.ExternalID,
			Title:       rawItem.Title,
			Content:     rawItem.Content,
			URL:         rawItem.URL,
			Author:      rawItem.Author,
			Status:      models.StatusUnread,
			Tags:        types.JSONArray[string]{},
			CategoryIDs: source.CategoryIDs,
			IsProcessed: false,
		}

		publishedAt := types.DateTime{}
		publishedAt.Scan(rawItem.PublishedAt)
		feedItem.PublishedAt = publishedAt

		fetchedAt := types.DateTime{}
		fetchedAt.Scan(time.Now())
		feedItem.FetchedAt = fetchedAt

		if len(rawItem.Tags) > 0 {
			tags := types.JSONArray[string]{}
			tags.Scan(rawItem.Tags)
			feedItem.Tags = tags
		}

		if rawItem.Metadata != nil {
			metadataJSON, err := json.Marshal(rawItem.Metadata)
			if err != nil {
				logger.LogError(fmt.Sprintf("Error marshaling metadata: %v", err))
			} else {
				feedItem.Metadata = string(metadataJSON)
			}
		}

		feedItem.Id = util.GenerateRandomId()
		// TODO: Should use upsert or insert
		if err := query.SaveRecord(feedItem); err != nil {
			logger.LogError(fmt.Sprintf("Error saving feed item: %v", err))
			continue
		}
	}

	now := types.DateTime{}
	now.Scan(time.Now())

	return query.UpdateRecord[*models.FeedSource](source.Id, map[string]interface{}{
		"last_fetched": now,
		"error_count":  0,
		"last_error":   "",
	})
}

// FetchAllSources fetches items from all active sources
func (s *FeedServiceImpl) FetchAllSources(ctx context.Context) error {
	sources, err := query.FindAllByFilter[*models.FeedSource](map[string]interface{}{
		"is_active": true,
	})

	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	errors := make(chan error, len(sources))
	semaphore := make(chan struct{}, 5) // Limit concurrent fetches

	for _, source := range sources {
		// Skip if refresh rate is set and not enough time has passed
		if source.RefreshRate > 0 {
			lastFetchTime := source.LastFetched.Time()
			nextFetchTime := lastFetchTime.Add(time.Duration(source.RefreshRate) * time.Minute)
			if time.Now().Before(nextFetchTime) {
				continue
			}
		}

		wg.Add(1)
		go func(src *models.FeedSource) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if err := s.FetchFromSource(ctx, src); err != nil {
				select {
				case errors <- err:
				default:
				}
			}
		}(source)
	}

	wg.Wait()
	close(errors)

	// Collect errors
	var errs []error
	for err := range errors {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("encountered %d errors while fetching sources", len(errs))
	}

	return nil
}

// ProcessNewItems processes all unprocessed items
func (s *FeedServiceImpl) ProcessNewItems(ctx context.Context) error {
	items, err := query.FindAllByFilter[*models.FeedItem](map[string]interface{}{
		"is_processed": false,
	})

	if err != nil {
		return err
	}

	if len(items) == 0 {
		logger.LogInfo("No new items to process")
		return nil
	}

	logger.LogInfo(fmt.Sprintf("Processing %d new feed items", len(items)))

	var wg sync.WaitGroup
	errors := make(chan error, len(items))
	semaphore := make(chan struct{}, 3) // Limit concurrent processing

	for _, item := range items {
		wg.Add(1)
		go func(itm *models.FeedItem) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if err := s.processor.ProcessItem(ctx, itm); err != nil {
				select {
				case errors <- err:
				default:
				}
				return
			}

			// Mark as processed
			query.UpdateRecord[*models.FeedItem](itm.Id, map[string]interface{}{
				"is_processed": true,
			})
		}(item)
	}

	wg.Wait()
	close(errors)

	// Collect errors
	var errs []error
	for err := range errors {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("encountered %d errors while processing items", len(errs))
	}

	return nil
}

// GetUserFeeds gets feed items for a user with filtering options
func (s *FeedServiceImpl) GetUserFeeds(ctx context.Context, userID string, options map[string]interface{}) ([]*models.FeedItem, error) {
	filter := map[string]interface{}{
		"user": userID,
	}

	// Add additional filters based on options
	for key, value := range options {
		switch key {
		case "status":
			filter["status"] = value
		case "category_id":
			// This is more complex as we need to filter by array membership
			// This may require custom SQL or a specific implementation in your query package
		case "source_id":
			filter["source_id"] = value
		case "search":
			// Text search is also more complex and may require custom implementation
		case "tags":
			// Array membership query
		case "limit", "offset", "sort":
			// These are pagination/sorting options, not filters
		}
	}

	// Get the items
	return query.FindAllByFilter[*models.FeedItem](filter)
}
