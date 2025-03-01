package ai

import (
	"context"

	"github.com/shashank-sharma/backend/internal/models"
)

// Request/response types for AI operations

// SummarizeRequest contains parameters for the summarization operation
type SummarizeRequest struct {
	Text      string `json:"text"`      // The text to summarize
	MaxLength int    `json:"maxLength"` // The maximum length of the summary in characters
}

// SummarizeResponse contains the result of a summarization operation
type SummarizeResponse struct {
	Summary string `json:"summary"` // The generated summary
}

// TagRequest contains parameters for content tagging
type TagRequest struct {
	Title   string `json:"title"`   // The title of the content
	Content string `json:"content"` // The content to analyze
	MaxTags int    `json:"maxTags"` // Maximum number of tags to generate
	UserID  string `json:"userId"`  // The user ID to associate tags with
}

// TagResponse contains the result of a tagging operation
type TagResponse struct {
	Tags     []string        `json:"tags"`     // The generated tag names (for backward compatibility)
	TagInfos []*models.Tag   `json:"tagInfos"` // The generated tag objects with full information
	TagIDs   []string        `json:"tagIds"`   // The IDs of the generated or found tags
}

// ClassifyRequest contains parameters for content classification
type ClassifyRequest struct {
	Title   string   `json:"title"`   // The title of the content
	Content string   `json:"content"` // The content to classify
	Labels  []string `json:"labels"`  // The possible classification labels
}

// ClassifyResponse contains the result of a classification operation
type ClassifyResponse struct {
	Label       string  `json:"label"`       // The primary classification label
	Confidence  float64 `json:"confidence"`  // Confidence score for the classification
	OtherLabels map[string]float64 `json:"otherLabels"` // Other possible labels with scores
}

// RecommendRequest contains parameters for making content recommendations
type RecommendRequest struct {
	UserID       string                 `json:"userId"`       // The user to make recommendations for
	Item         *models.FeedItem       `json:"item"`         // The item to compare against
	UserMetadata map[string]interface{} `json:"userMetadata"` // User preferences and history
}

// RecommendResponse contains the result of a recommendation operation
type RecommendResponse struct {
	Score       float64 `json:"score"`       // Relevance score (0-1)
	Explanation string  `json:"explanation"` // Explanation of the recommendation
}

// AIClient defines the interface for AI services
type AIClient interface {
	// Summarize generates a concise summary of the provided text
	Summarize(ctx context.Context, req *SummarizeRequest) (*SummarizeResponse, error)

	// SuggestTags extracts relevant tags from the content
	SuggestTags(ctx context.Context, req *TagRequest) (*TagResponse, error)
	
	// ClassifyContent classifies content into predefined categories
	ClassifyContent(ctx context.Context, req *ClassifyRequest) (*ClassifyResponse, error)
	
	// RecommendContent provides a relevance score for content based on user preferences
	RecommendContent(ctx context.Context, req *RecommendRequest) (*RecommendResponse, error)
} 