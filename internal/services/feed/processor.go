package feed

import (
	"context"
	"fmt"
	"strings"

	"github.com/shashank-sharma/backend/internal/logger"
	"github.com/shashank-sharma/backend/internal/models"
	"github.com/shashank-sharma/backend/internal/query"
	"github.com/shashank-sharma/backend/internal/services"
	"github.com/shashank-sharma/backend/internal/services/ai"
	"github.com/shashank-sharma/backend/internal/util"
)

type FeedProcessorImpl struct {
	aiClient ai.AIClient
}

func NewFeedProcessor(aiClient ai.AIClient) services.FeedProcessor {
	return &FeedProcessorImpl{
		aiClient: aiClient,
	}
}

// ProcessItem processes a feed item
func (p *FeedProcessorImpl) ProcessItem(ctx context.Context, item *models.FeedItem) error {
	logger.LogInfo(fmt.Sprintf("Processing feed item: %s", item.Title))

	// 1. Generate summary if content is available
	if item.Content != "" {
		summary, err := p.GenerateSummary(ctx, item)
		if err != nil {
			logger.LogError(fmt.Sprintf("Error generating summary: %v", err))
		} else {
			item.Summary = summary
		}
	}

	// 2. Suggest tags if not already present
	if len(item.Tags) == 0 {
		tags, err := p.SuggestTags(ctx, item)
		if err != nil {
			logger.LogError(fmt.Sprintf("Error suggesting tags: %v", err))
		} else {
			for _, tag := range tags {
				item.Tags = append(item.Tags, tag)
			}
		}
	}

	// 3. Update the processed item in the database
	return query.UpdateRecord[*models.FeedItem](item.Id, map[string]interface{}{
		"summary":      item.Summary,
		"tags":         item.Tags,
		"is_processed": true,
	})
}

// GenerateSummary generates a summary for the feed item using AI
func (p *FeedProcessorImpl) GenerateSummary(ctx context.Context, item *models.FeedItem) (string, error) {
	// If AI client isn't configured, use a simple fallback approach
	// Maybe NLP ?
	if p.aiClient == nil {
		return p.fallbackGenerateSummary(item)
	}

	// Extract clean text from content
	content := stripTags(item.Content)
	if content == "" {
		return "", nil
	}

	// Use AI client to summarize
	resp, err := p.aiClient.Summarize(ctx, &ai.SummarizeRequest{
		Text:      content,
		MaxLength: 150, // Default summary length
	})

	if err != nil {
		logger.LogError(fmt.Sprintf("AI summarization error: %v", err))
		// If AI fails, fall back to simple approach
		return p.fallbackGenerateSummary(item)
	}

	return resp.Summary, nil
}

// fallbackGenerateSummary provides a simple non-AI summary when AI is unavailable
func (p *FeedProcessorImpl) fallbackGenerateSummary(item *models.FeedItem) (string, error) {
	content := stripTags(item.Content)
	if len(content) > 150 {
		return content[:150] + "...", nil
	}
	return content, nil
}

// SuggestTags suggests tags for the feed item using AI and returns tag IDs
func (p *FeedProcessorImpl) SuggestTags(ctx context.Context, item *models.FeedItem) ([]string, error) {
	var tagNames []string
	var err error
	
	if p.aiClient == nil {
		tagNames, err = p.fallbackSuggestTagNames(item)
	} else {
		// Extract clean text from content
		content := stripTags(item.Content)
		
		// Use AI client to suggest tags
		resp, err := p.aiClient.SuggestTags(ctx, &ai.TagRequest{
			Title:   item.Title,
			Content: content,
			MaxTags: 5, // Default max tags
		})

		if err != nil {
			logger.LogError(fmt.Sprintf("AI tagging error: %v", err))
			tagNames, err = p.fallbackSuggestTagNames(item)
		} else {
			tagNames = resp.Tags
		}
	}
	
	if err != nil {
		return nil, err
	}
	
	var tagIDs []string
	
	for _, tagName := range tagNames {
		if tagName == "" {
			continue
		}
		
		// Normalize tag name (lowercase, trim)
		tagName = strings.ToLower(strings.TrimSpace(tagName))
		
		// Check if tag already exists for this user
		existingTag, err := query.FindByFilter[*models.Tag](map[string]interface{}{
			"user": item.User,
			"name": tagName,
		})
		
		if err == nil && existingTag != nil {
			// Tag exists, use its ID
			tagIDs = append(tagIDs, existingTag.Id)
		} else {
			// Tag doesn't exist, create a new one
			// TODO: Use query.Upsert
			newTag := &models.Tag{
				User:        item.User,
				Name:        tagName,
				Color:       generateTagColor(tagName),
				Description: "", // Empty description for AI-generated tags
				IsAICreated: true,
			}
			
			newTag.Id = util.GenerateRandomId()
			if err := query.SaveRecord(newTag); err != nil {
				logger.LogError(fmt.Sprintf("Error saving new tag: %v", err))
				continue
			}
			
			tagIDs = append(tagIDs, newTag.Id)
		}
	}
	
	return tagIDs, nil
}

// fallbackSuggestTagNames provides simple keyword-based tagging when AI is unavailable
func (p *FeedProcessorImpl) fallbackSuggestTagNames(item *models.FeedItem) ([]string, error) {
	// Combine title and content for analysis
	text := item.Title + " " + stripTags(item.Content)
	text = strings.ToLower(text)

	// Common tech keywords - basic keyword matching
	keywords := []string{
		"golang", "python", "javascript", "react", "svelte", "web", "mobile",
		"database", "cloud", "api", "security", "design", "blockchain", "ai",
		"machine learning", "data science", "devops", "backend", "frontend",
	}

	var tags []string
	for _, keyword := range keywords {
		if strings.Contains(text, keyword) {
			tags = append(tags, keyword)
		}

		// Limit to 5 tags
		if len(tags) >= 5 {
			break
		}
	}

	return tags, nil
}

// Helper function to generate a consistent color based on the tag name
func generateTagColor(tagName string) string {
	// Simple hash function to generate a color
	var hash int
	for _, c := range tagName {
		hash = (hash*31 + int(c)) % 360
	}
	
	// Generate HSL color with good saturation and lightness
	return fmt.Sprintf("hsl(%d, 70%%, 60%%)", hash)
}

// Helper function to strip HTML tags from text
func stripTags(html string) string {
	// This is a very simplistic implementation
	// TODO: use a proper HTML parser
	stripped := strings.ReplaceAll(html, "<br>", "\n")
	stripped = strings.ReplaceAll(stripped, "<p>", "\n")
	stripped = strings.ReplaceAll(stripped, "</p>", "\n")

	var result strings.Builder
	var tagOpen bool

	for _, r := range stripped {
		if r == '<' {
			tagOpen = true
			continue
		}
		if r == '>' {
			tagOpen = false
			continue
		}
		if !tagOpen {
			result.WriteRune(r)
		}
	}

	return result.String()
}
