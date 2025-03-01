package ai

import (
	"fmt"

	"github.com/shashank-sharma/backend/internal/config"
	"github.com/shashank-sharma/backend/internal/logger"
)

// NewAIClient creates an AI client based on the provided configuration
func NewAIClient(cfg config.AIConfig) (AIClient, error) {
	switch cfg.Service {
	case config.AIServiceOpenAI:
		if cfg.APIKey == "" {
			return nil, fmt.Errorf("OpenAI API key is required")
		}
		logger.LogInfo(fmt.Sprintf("Initializing OpenAI client with model: %s", cfg.Model))
		return NewOpenAIClient(cfg.APIKey, cfg.Model), nil
	
	case config.AIServiceClaude:
		if cfg.AnthropicKey == "" {
			return nil, fmt.Errorf("Anthropic API key is required for Claude")
		}
		logger.LogInfo(fmt.Sprintf("Initializing Claude client with model: %s", cfg.Model))
		return NewClaudeClient(cfg.AnthropicKey, cfg.Model), nil
		
	case config.AIServiceNone:
		logger.LogInfo("AI services disabled")
		return nil, nil
		
	default:
		return nil, fmt.Errorf("unsupported AI service: %s", cfg.Service)
	}
} 