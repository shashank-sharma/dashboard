package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/shashank-sharma/backend/internal/logger"
	"github.com/shashank-sharma/backend/internal/models"
	"github.com/shashank-sharma/backend/internal/services"
)

const (
	HNItemEndpoint   = "https://hacker-news.firebaseio.com/v0/item/%d.json"
	HNTopStoriesURL  = "https://hacker-news.firebaseio.com/v0/topstories.json"
	HNNewStoriesURL  = "https://hacker-news.firebaseio.com/v0/newstories.json"
	HNBestStoriesURL = "https://hacker-news.firebaseio.com/v0/beststories.json"
	MaxConcurrentFetches = 5
)

// HackerNewsProvider implements the FeedSourceProvider interface for Hacker News
type HackerNewsProvider struct {
	client         *http.Client
	contentFetcher *HTMLContentFetcher
	fetchContent   bool
}

type HNItem struct {
	ID          int    `json:"id"`
	Type        string `json:"type"`
	By          string `json:"by"`
	Time        int64  `json:"time"`
	Text        string `json:"text"`
	URL         string `json:"url"`
	Title       string `json:"title"`
	Score       int    `json:"score"`
	Kids        []int  `json:"kids"`
	Descendants int    `json:"descendants"`
}

func NewHackerNewsProvider() services.FeedSourceProvider {
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}
	
	return &HackerNewsProvider{
		client:         httpClient,
		contentFetcher: NewHTMLContentFetcher(60 * time.Second), // Longer timeout for content fetching
		fetchContent:   true, // Enable content fetching by default
	}
}

// GetProviderType returns the type of feed source
func (p *HackerNewsProvider) GetProviderType() string {
	return models.SourceTypeHackerNews
}

// FetchItems fetches items from Hacker News
func (p *HackerNewsProvider) FetchItems(ctx context.Context, config map[string]interface{}, lastFetched time.Time) ([]services.RawFeedItem, error) {
	// Get feed type (top, new, best)
	feedType, ok := config["feed_type"].(string)
	if !ok || feedType == "" {
		feedType = "top" // Default to top stories
	}

	limit := 30 // Default limit
	if limitVal, ok := config["limit"].(float64); ok {
		limit = int(limitVal)
	}
	
	if fetchContent, ok := config["fetch_content"].(bool); ok {
		p.fetchContent = fetchContent
	}

	// Get the story IDs URL based on feed type
	var storiesURL string
	switch feedType {
	case "top":
		storiesURL = HNTopStoriesURL
	case "new":
		storiesURL = HNNewStoriesURL
	case "best":
		storiesURL = HNBestStoriesURL
	default:
		return nil, fmt.Errorf("invalid feed type: %s", feedType)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", storiesURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching story IDs: %v", err)
	}
	defer resp.Body.Close()

	var storyIDs []int
	if err := json.NewDecoder(resp.Body).Decode(&storyIDs); err != nil {
		return nil, fmt.Errorf("error decoding story IDs: %v", err)
	}

	// Limit the number of stories to fetch
	if len(storyIDs) > limit {
		storyIDs = storyIDs[:limit]
	}

	// Fetch each story
	var items []services.RawFeedItem
	
	// Create a semaphore to limit concurrent fetches
	semaphore := make(chan struct{}, MaxConcurrentFetches)
	itemsChan := make(chan services.RawFeedItem)
	doneChan := make(chan bool)
	finishedChan := make(chan struct{})
	
	activeJobs := 0
	
	go func() {
		for item := range itemsChan {
			items = append(items, item)
		}
		finishedChan <- struct{}{}
	}()
	
	ctxDone := false
	for _, id := range storyIDs {
		select {
		case <-ctx.Done():
			ctxDone = true
			break
		default:
			// Continue processing
		}
		
		if ctxDone {
			break
		}
		
		semaphore <- struct{}{}
		activeJobs++
		
		go func(storyID int) {
			defer func() {
				<-semaphore
				doneChan <- true
			}()
			
			item, err := p.fetchItem(ctx, storyID)
			if err != nil {
				logger.LogError(fmt.Sprintf("Error fetching Hacker News item %d: %v", storyID, err))
				return
			}

			itemTime := time.Unix(item.Time, 0)
			if !lastFetched.IsZero() && !itemTime.After(lastFetched) {
				return
			}

			// Only include stories, not comments
			if item.Type != "story" || item.Title == "" {
				return
			}
			
			// Use the text from HN if available, otherwise fetch content if URL is present
			content := item.Text
			
			// If no content but there's a URL and content fetching is enabled, try to fetch content
			// TODO: Support content fetch for PDF/other file formats
			if content == "" && item.URL != "" && p.fetchContent {
				fetchedContent, err := p.contentFetcher.FetchContent(ctx, item.URL)
				if err != nil {
					logger.LogError(fmt.Sprintf("Error fetching content for %s: %v", item.URL, err))
				} else if fetchedContent != "" {
					content = fetchedContent
					logger.LogInfo(fmt.Sprintf("Successfully fetched content for HN item %d: %s", storyID, item.URL))
				}
			}

			// Create feed item
			feedItem := services.RawFeedItem{
				ExternalID:  strconv.Itoa(item.ID),
				Title:       item.Title,
				Content:     content,
				URL:         item.URL,
				Author:      item.By,
				PublishedAt: itemTime,
				Tags:        []string{"hackernews", feedType},
				Metadata: map[string]interface{}{
					"score":         item.Score,
					"comment_count": item.Descendants,
					"type":          item.Type,
					"hn_url":        fmt.Sprintf("https://news.ycombinator.com/item?id=%d", item.ID),
					"content_fetched": content != "" && content != item.Text,
				},
			}

			select {
			case itemsChan <- feedItem:
				// Item sent successfully
			case <-ctx.Done():
				// Context was canceled, don't bother sending
			}
		}(id)
	}
	
	completedJobs := 0
	for completedJobs < activeJobs {
		select {
		case <-doneChan:
			completedJobs++
		case <-ctx.Done():
			logger.LogError("Context canceled while fetching Hacker News items")
			ctxDone = true
		}
	}
	
	close(itemsChan)
	<-finishedChan

	logger.LogInfo(fmt.Sprintf("Fetched %d items from Hacker News", len(items)))
	return items, nil
}

// fetchItem fetches a single item from the Hacker News API
func (p *HackerNewsProvider) fetchItem(ctx context.Context, id int) (*HNItem, error) {
	url := fmt.Sprintf(HNItemEndpoint, id)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var item HNItem
	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		return nil, err
	}

	return &item, nil
}

// Validate validates the source configuration
func (p *HackerNewsProvider) Validate(config map[string]interface{}) error {
	if feedType, ok := config["feed_type"].(string); ok {
		switch feedType {
		case "top", "new", "best":
			// Valid feed types
		default:
			return fmt.Errorf("invalid feed type: %s", feedType)
		}
	}

	if limit, ok := config["limit"].(float64); ok {
		if limit < 1 || limit > 100 {
			return fmt.Errorf("limit must be between 1 and 100")
		}
	}
	
	if fetchContent, ok := config["fetch_content"].(bool); ok {
		_ = fetchContent
	}

	return nil
}
