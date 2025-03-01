package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/shashank-sharma/backend/internal/logger"
	"github.com/shashank-sharma/backend/internal/models"
	"github.com/shashank-sharma/backend/internal/query"
	"github.com/shashank-sharma/backend/internal/services"
	"github.com/shashank-sharma/backend/internal/util"
)

type CreateFeedSourceRequest struct {
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	URL         string                 `json:"url"`
	Config      map[string]interface{} `json:"config"`
	CategoryIDs []string               `json:"category_ids"`
	RefreshRate int                    `json:"refresh_rate"`
}

type FeedSourceResponse struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	URL         string                 `json:"url"`
	Config      map[string]interface{} `json:"config"`
	CategoryIDs []string               `json:"category_ids"`
	RefreshRate int                    `json:"refresh_rate"`
	IsActive    bool                   `json:"is_active"`
	LastFetched string                 `json:"last_fetched"`
	ErrorCount  int                    `json:"error_count"`
	LastError   string                 `json:"last_error"`
}

type FeedCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Icon        string `json:"icon"`
	Type        string `json:"type"`
	SortOrder   int    `json:"sort_order"`
	IsDefault   bool   `json:"is_default"`
}

type FeedItemResponse struct {
	ID            string                 `json:"id"`
	SourceID      string                 `json:"source_id"`
	SourceName    string                 `json:"source_name"`
	Title         string                 `json:"title"`
	URL           string                 `json:"url"`
	Author        string                 `json:"author"`
	Summary       string                 `json:"summary"`
	PublishedAt   string                 `json:"published_at"`
	Status        string                 `json:"status"`
	Rating        int                    `json:"rating"`
	TagIDs        []string               `json:"tag_ids"`
	Tags          []TagResponse          `json:"tags"`
	CategoryIDs   []string               `json:"category_ids"`
	CategoryNames []string               `json:"category_names"`
	Metadata      map[string]interface{} `json:"metadata"`
}

type TagResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Color       string `json:"color"`
	IsAICreated bool   `json:"is_ai_created"`
}

// RegisterFeedRoutes registers feed-related API routes
func RegisterFeedRoutes(e *core.ServeEvent, feedService services.FeedService) {
	e.Router.POST("/api/feeds/sources", func(e *core.RequestEvent) error {
		return CreateFeedSource(e, feedService)
	})
	e.Router.POST("/api/feeds/sources/{id}/fetch", func(e *core.RequestEvent) error {
		return FetchFromSource(e, feedService)
	})
	e.Router.GET("/api/feeds", func(e *core.RequestEvent) error {
		return GetFeeds(e, feedService)
	})
}

// CreateFeedSource creates a new feed source
func CreateFeedSource(e *core.RequestEvent, feedService services.FeedService) error {
	// Get user ID from token
	token := e.Request.Header.Get("Authorization")
	userId, err := util.GetUserId(token)
	if err != nil {
		return e.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Unauthorized"})
	}

	// Parse request body
	req := &CreateFeedSourceRequest{}
	if err := e.BindBody(req); err != nil {
		return e.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid request body"})
	}

	// Validate source type
	provider, exists := feedService.GetProvider(req.Type)
	if !exists {
		return e.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Unsupported feed source type",
		})
	}

	// Validate provider-specific configuration
	if err := provider.Validate(req.Config); err != nil {
		return e.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Convert config to JSON string
	configJSON, err := json.Marshal(req.Config)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid configuration format",
		})
	}

	// Create feed source
	categoryIDs := types.JSONArray[string]{}
	if len(req.CategoryIDs) > 0 {
		categoryIDs.Scan(req.CategoryIDs)
	}

	source := &models.FeedSource{
		User:        userId,
		Name:        req.Name,
		Type:        req.Type,
		URL:         req.URL,
		Config:      string(configJSON),
		CategoryIDs: categoryIDs,
		RefreshRate: req.RefreshRate,
		IsActive:    true,
	}

	// Save to database
	source.Id = util.GenerateRandomId()
	if err := query.SaveRecord(source); err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to create feed source",
		})
	}

	return e.JSON(http.StatusCreated, map[string]interface{}{
		"id":      source.Id,
		"message": "Feed source created successfully",
	})
}

// FetchFromSource manually fetches from a source
func FetchFromSource(e *core.RequestEvent, feedService services.FeedService) error {
	// Get user ID from token
	token := e.Request.Header.Get("Authorization")
	userId, err := util.GetUserId(token)
	if err != nil {
		return e.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Unauthorized"})
	}

	// Get source ID from URL
	sourceId := e.Request.PathValue("id")
	if sourceId == "" {
		return e.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Missing source ID"})
	}

	// Get source from database
	source, err := query.FindByFilter[*models.FeedSource](map[string]interface{}{
		"id":   sourceId,
		"user": userId,
	})
	if err != nil {
		return e.JSON(http.StatusNotFound, map[string]interface{}{"error": "Feed source not found"})
	}

	// Fetch from source - create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Fetch from source
	if err := feedService.FetchFromSource(ctx, source); err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to fetch from source: " + err.Error(),
		})
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"message": "Successfully fetched from source",
	})
}

// GetFeeds returns a list of feed items for the authenticated user
func GetFeeds(e *core.RequestEvent, feedService services.FeedService) error {
	// Get authenticated user
	token := e.Request.Header.Get("Authorization")
	userId, err := util.GetUserId(token)
	if err != nil {
		return e.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Unauthorized"})
	}

	// Parse query parameters
	var parseErr error
	source := e.Request.URL.Query().Get("source")
	status := e.Request.URL.Query().Get("status")
	limit := 50
	if limitStr := e.Request.URL.Query().Get("limit"); limitStr != "" {
		limit, parseErr = strconv.Atoi(limitStr)
		if parseErr != nil {
			limit = 50
		}
	}
	offset := 0
	if offsetStr := e.Request.URL.Query().Get("offset"); offsetStr != "" {
		offset, parseErr = strconv.Atoi(offsetStr)
		if parseErr != nil {
			offset = 0
		}
	}

	// Build query filter
	filter := map[string]interface{}{
		"user": userId,
	}

	if source != "" {
		filter["source_id"] = source
	}
	if status != "" {
		filter["status"] = status
	}

	// Fetch items
	items, err := query.FindAllByFilterWithPagination[*models.FeedItem](filter, limit, offset)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to fetch feed items: " + err.Error(),
		})
	}

	sources, err := query.FindAllByFilter[*models.FeedSource](map[string]interface{}{
		"user": userId,
	})
	if err != nil {
		logger.LogError(err.Error())
	}
	sourceMap := make(map[string]string)
	for _, source := range sources {
		sourceMap[source.Id] = source.Name
	}

	categories, err := query.FindAllByFilter[*models.FeedCategory](map[string]interface{}{
		"user": userId,
	})
	if err != nil {
		logger.LogError(err.Error())
	}
	categoryMap := make(map[string]string)
	for _, category := range categories {
		categoryMap[category.Id] = category.Name
	}

	tags, err := query.FindAllByFilter[*models.Tag](map[string]interface{}{
		"user": userId,
	})
	if err != nil {
		logger.LogError(err.Error())
	}

	tagMap := make(map[string]*models.Tag)
	for _, tag := range tags {
		tagMap[tag.Id] = tag
	}

	response := make([]FeedItemResponse, 0, len(items))
	for _, item := range items {
		sourceName := ""
		if name, ok := sourceMap[item.SourceID]; ok {
			sourceName = name
		}

		categoryNames := make([]string, 0)
		for _, catID := range item.CategoryIDs {
			if name, ok := categoryMap[catID]; ok {
				categoryNames = append(categoryNames, name)
			}
		}

		metadata := make(map[string]interface{})
		if item.Metadata != "" {
			if err := json.Unmarshal([]byte(item.Metadata), &metadata); err != nil {
				logger.LogError("Failed to parse metadata: " + err.Error())
			}
		}

		tagResponses := make([]TagResponse, 0)
		for _, tagID := range item.Tags {
			if tag, ok := tagMap[tagID]; ok {
				tagResponses = append(tagResponses, TagResponse{
					ID:          tag.Id,
					Name:        tag.Name,
					Color:       tag.Color,
					IsAICreated: tag.IsAICreated,
				})
			}
		}

		response = append(response, FeedItemResponse{
			ID:            item.Id,
			SourceID:      item.SourceID,
			SourceName:    sourceName,
			Title:         item.Title,
			URL:           item.URL,
			Author:        item.Author,
			Summary:       item.Summary,
			PublishedAt:   item.PublishedAt.String(),
			Status:        item.Status,
			Rating:        item.Rating,
			TagIDs:        item.Tags,
			Tags:          tagResponses,
			CategoryIDs:   item.CategoryIDs,
			CategoryNames: categoryNames,
			Metadata:      metadata,
		})
	}

	return e.JSON(http.StatusOK, response)
}