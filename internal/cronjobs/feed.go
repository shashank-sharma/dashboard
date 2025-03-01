package cronjobs

import (
	"context"
	"fmt"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/shashank-sharma/backend/internal/logger"
	"github.com/shashank-sharma/backend/internal/services"
)

// FeedUpdateJob runs scheduled feed updates
func FeedUpdateJob(app *pocketbase.PocketBase, feedService services.FeedService) error {
	logger.LogInfo("Starting feed update job")

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Fetch from all sources
	if err := feedService.FetchAllSources(ctx); err != nil {
		logger.LogError(fmt.Sprintf("Error fetching feeds: %v", err))
		return err
	}

	// Process new items
	if err := feedService.ProcessNewItems(ctx); err != nil {
		logger.LogError(fmt.Sprintf("Error processing feed items: %v", err))
		return err
	}

	logger.LogInfo("Feed update job completed")
	return nil
}
