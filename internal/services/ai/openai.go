package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"
	"github.com/shashank-sharma/backend/internal/logger"
	"github.com/shashank-sharma/backend/internal/models"
)

// OpenAIClient implements the AIClient interface using OpenAI's API
type OpenAIClient struct {
	client *openai.Client
	model  string
}

// NewOpenAIClient creates a new OpenAI client with the provided API key
func NewOpenAIClient(apiKey string, model string) AIClient {
	client := openai.NewClient(apiKey)
	
	if model == "" {
		model = openai.GPT3Dot5Turbo
	}
	
	return &OpenAIClient{
		client: client,
		model:  model,
	}
}

// Summarize implements the AIClient.Summarize method
func (c *OpenAIClient) Summarize(ctx context.Context, req *SummarizeRequest) (*SummarizeResponse, error) {
	if req.Text == "" {
		return &SummarizeResponse{Summary: ""}, nil
	}
	
	// Default max length
	maxLength := 150
	if req.MaxLength > 0 {
		maxLength = req.MaxLength
	}
	
	// Create the prompt
	prompt := fmt.Sprintf(
		"Summarize the following text in a concise, informative way in %d characters or less:\n\n%s",
		maxLength,
		req.Text,
	)
	
	// Create the chat completion request
	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: c.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens:   int(float64(maxLength) * 0.5), // Estimate tokens from characters
			Temperature: 0.3, // Lower temperature for more focused summaries
		},
	)
	
	if err != nil {
		logger.LogError(fmt.Sprintf("OpenAI summarization error: %v", err))
		return nil, fmt.Errorf("failed to generate summary: %w", err)
	}
	
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no summary generated")
	}
	
	summary := strings.TrimSpace(resp.Choices[0].Message.Content)
	
	return &SummarizeResponse{
		Summary: summary,
	}, nil
}

// SuggestTags implements the AIClient.SuggestTags method
func (c *OpenAIClient) SuggestTags(ctx context.Context, req *TagRequest) (*TagResponse, error) {
	if req.Content == "" && req.Title == "" {
		return &TagResponse{Tags: []string{}, TagInfos: []*models.Tag{}, TagIDs: []string{}}, nil
	}
	
	// Default max tags
	maxTags := 5
	if req.MaxTags > 0 {
		maxTags = req.MaxTags
	}
	
	// Create the prompt
	content := req.Content
	if content == "" {
		content = req.Title
	} else if req.Title != "" {
		content = req.Title + "\n\n" + content
	}
	
	prompt := fmt.Sprintf(
		"Extract up to %d relevant tags from the following content. Return only a JSON array of tag strings, with no explanations:\n\n%s",
		maxTags,
		content,
	)
	
	// Create the chat completion request
	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: c.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a tagging assistant. Your job is to extract relevant tags from content. Always return tags as a JSON array of strings.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens:   100,
			Temperature: 0.3,
		},
	)
	
	if err != nil {
		logger.LogError(fmt.Sprintf("OpenAI tagging error: %v", err))
		return nil, fmt.Errorf("failed to generate tags: %w", err)
	}
	
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no tags generated")
	}
	
	// Parse the JSON response
	tagsContent := strings.TrimSpace(resp.Choices[0].Message.Content)
	
	// Handle case where the AI might include explanation text
	if !strings.HasPrefix(tagsContent, "[") {
		// Try to find JSON array in the response
		startIdx := strings.Index(tagsContent, "[")
		endIdx := strings.LastIndex(tagsContent, "]")
		if startIdx >= 0 && endIdx > startIdx {
			tagsContent = tagsContent[startIdx : endIdx+1]
		} else {
			// If no JSON array found, try to extract words that look like tags
			words := strings.Fields(tagsContent)
			var extractedTags []string
			for _, word := range words {
				// Clean the word
				word = strings.Trim(word, `",.;:()[]{}`)
				if len(word) > 2 { // Only consider words of reasonable length
					extractedTags = append(extractedTags, word)
				}
			}
			if len(extractedTags) > 0 {
				// Convert the extracted words to a JSON array string
				tagBytes, _ := json.Marshal(extractedTags)
				tagsContent = string(tagBytes)
			} else {
				return nil, fmt.Errorf("could not parse tags from response: %s", tagsContent)
			}
		}
	}
	
	// Parse the JSON array
	var tagNames []string
	if err := json.Unmarshal([]byte(tagsContent), &tagNames); err != nil {
		return nil, fmt.Errorf("failed to parse tags JSON: %w", err)
	}
	
	// Limit the number of tags
	if len(tagNames) > maxTags {
		tagNames = tagNames[:maxTags]
	}
	
	// Create Tag models for each tag name
	tagInfos := make([]*models.Tag, 0, len(tagNames))
	tagIDs := make([]string, 0, len(tagNames))
	
	for _, name := range tagNames {
		// Clean the tag name
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		
		// Convert tag name to a consistent color based on the name
		// This ensures the same tag always gets the same color
		colorHex := generateTagColor(name)
		
		// Create a Tag model
		tag := &models.Tag{
			User:        req.UserID,
			Name:        name,
			Color:       colorHex,
			Description: fmt.Sprintf("AI-generated tag for content related to %s", name),
			IsAICreated: true,
		}
		
		// Add a random ID for now - actual saving and ID generation will happen in the processor
		tagInfos = append(tagInfos, tag)
	}
	
	return &TagResponse{
		Tags:     tagNames,
		TagInfos: tagInfos,
		TagIDs:   tagIDs, // Will be populated when the tags are saved to the database
	}, nil
}

// generateTagColor creates a deterministic color based on the tag name
func generateTagColor(tagName string) string {
	// Simple hash function to get a deterministic number from a string
	var hash uint32
	for i := 0; i < len(tagName); i++ {
		hash = hash*31 + uint32(tagName[i])
	}
	
	// Use HSL color space for better distribution
	// We'll use the hash to determine hue (0-360) while keeping saturation and lightness fixed
	hue := hash % 360
	
	// Convert HSL to RGB
	// For simplicity, we're using a fixed saturation and lightness
	// S = 65%, L = 50% gives vibrant but not too bright colors
	r, g, b := hslToRgb(float64(hue), 0.65, 0.5)
	
	// Return as hex color code
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

// hslToRgb converts HSL color values to RGB
func hslToRgb(h, s, l float64) (r, g, b uint8) {
	h = h / 360.0 // Convert to 0-1 range
	
	var v1, v2 float64
	if s == 0 {
		// Achromatic (grey)
		r = uint8(l * 255)
		g = uint8(l * 255)
		b = uint8(l * 255)
		return
	}
	
	if l < 0.5 {
		v2 = l * (1 + s)
	} else {
		v2 = (l + s) - (s * l)
	}
	
	v1 = 2*l - v2
	
	r = uint8(255 * hueToRgb(v1, v2, h+1.0/3.0))
	g = uint8(255 * hueToRgb(v1, v2, h))
	b = uint8(255 * hueToRgb(v1, v2, h-1.0/3.0))
	return
}

// hueToRgb is a helper function for hslToRgb
func hueToRgb(v1, v2, h float64) float64 {
	if h < 0 {
		h += 1
	}
	if h > 1 {
		h -= 1
	}
	if 6*h < 1 {
		return v1 + (v2-v1)*6*h
	}
	if 2*h < 1 {
		return v2
	}
	if 3*h < 2 {
		return v1 + (v2-v1)*(2.0/3.0-h)*6
	}
	return v1
}

// ClassifyContent implements the AIClient.ClassifyContent method
func (c *OpenAIClient) ClassifyContent(ctx context.Context, req *ClassifyRequest) (*ClassifyResponse, error) {
	if (req.Content == "" && req.Title == "") || len(req.Labels) == 0 {
		return &ClassifyResponse{
			Label:       "",
			Confidence:  0,
			OtherLabels: map[string]float64{},
		}, nil
	}
	
	// Create the prompt
	content := req.Content
	if content == "" {
		content = req.Title
	} else if req.Title != "" {
		content = req.Title + "\n\n" + content
	}
	
	labelsStr := strings.Join(req.Labels, ", ")
	prompt := fmt.Sprintf(
		"Classify the following content into one of these categories: %s.\nReturn a JSON object with keys: 'label' (string), 'confidence' (float 0-1), and 'otherLabels' (map of label to confidence).\n\nContent:\n%s",
		labelsStr,
		content,
	)
	
	// Create the chat completion request
	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: c.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a classification assistant. Your job is to classify content into predefined categories. Always return a JSON object with the structure: {\"label\": string, \"confidence\": float, \"otherLabels\": {string: float}}",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens:   150,
			Temperature: 0.2,
		},
	)
	
	if err != nil {
		logger.LogError(fmt.Sprintf("OpenAI classification error: %v", err))
		return nil, fmt.Errorf("failed to classify content: %w", err)
	}
	
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no classification generated")
	}
	
	// Parse the JSON response
	respContent := strings.TrimSpace(resp.Choices[0].Message.Content)
	
	// Extract JSON from the response
	startIdx := strings.Index(respContent, "{")
	endIdx := strings.LastIndex(respContent, "}")
	if startIdx >= 0 && endIdx > startIdx {
		respContent = respContent[startIdx : endIdx+1]
	}
	
	var result ClassifyResponse
	if err := json.Unmarshal([]byte(respContent), &result); err != nil {
		logger.LogError(fmt.Sprintf("Failed to parse classification JSON: %v, content: %s", err, respContent))
		
		// Fallback to first label with zero confidence
		if len(req.Labels) > 0 {
			return &ClassifyResponse{
				Label:       req.Labels[0],
				Confidence:  0,
				OtherLabels: map[string]float64{},
			}, nil
		}
		return nil, fmt.Errorf("failed to parse classification response: %w", err)
	}
	
	return &result, nil
}

// RecommendContent implements the AIClient.RecommendContent method
func (c *OpenAIClient) RecommendContent(ctx context.Context, req *RecommendRequest) (*RecommendResponse, error) {
	if req.Item == nil {
		return nil, fmt.Errorf("item cannot be nil")
	}
	
	// Build a user profile based on available metadata
	userProfileStr := "User preferences and history:\n"
	for k, v := range req.UserMetadata {
		userProfileStr += fmt.Sprintf("- %s: %v\n", k, v)
	}
	
	// Create item description
	itemDescStr := fmt.Sprintf("Content item:\n- Title: %s\n- URL: %s\n", req.Item.Title, req.Item.URL)
	
	// Add tags if available
	if len(req.Item.Tags) > 0 {
		itemDescStr += "- Tags: " + strings.Join(req.Item.Tags, ", ") + "\n"
	}
	
	// Add summary if available
	if req.Item.Summary != "" {
		itemDescStr += "- Summary: " + req.Item.Summary + "\n"
	}
	
	// Create the prompt
	prompt := fmt.Sprintf(
		"Based on the user profile, evaluate if the user would be interested in this content item. Score from 0 to 1, where 1 means highly relevant:\n\n%s\n\n%s\n\nProvide your response in the following JSON format only:\n{\"score\": 0.X, \"explanation\": \"Your reasoning here\"}",
		userProfileStr,
		itemDescStr,
	)
	
	// Create the chat completion request
	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: c.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a content recommendation engine. Analyze user preferences and content characteristics to predict relevance.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens:   150,
			Temperature: 0.3,
		},
	)
	
	if err != nil {
		logger.LogError(fmt.Sprintf("OpenAI recommendation error: %v", err))
		return nil, fmt.Errorf("failed to generate recommendation: %w", err)
	}
	
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no recommendation generated")
	}
	
	// Parse the JSON response
	respContent := strings.TrimSpace(resp.Choices[0].Message.Content)
	
	// Extract JSON from the response
	startIdx := strings.Index(respContent, "{")
	endIdx := strings.LastIndex(respContent, "}")
	if startIdx >= 0 && endIdx > startIdx {
		respContent = respContent[startIdx : endIdx+1]
	}
	
	var result RecommendResponse
	if err := json.Unmarshal([]byte(respContent), &result); err != nil {
		logger.LogError(fmt.Sprintf("Failed to parse recommendation JSON: %v, content: %s", err, respContent))
		
		// Fallback to neutral score
		return &RecommendResponse{
			Score:       0.5,
			Explanation: "Error parsing recommendation response",
		}, nil
	}
	
	return &result, nil
} 