package config

import (
	"os"
)

const (
	AIServiceOpenAI = "openai"
	AIServiceClaude = "claude"
	AIServiceNone   = "none"
	
	EnvAIService         = "AI_SERVICE"
	EnvAIAPIKey          = "AI_API_KEY"
	EnvAIModel           = "AI_MODEL"
	EnvAIAnthropicAPIKey = "AI_ANTHROPIC_API_KEY"
	
	DefaultAIService = AIServiceNone
	DefaultOpenAIModel = "gpt-3.5-turbo"
	DefaultClaudeModel = "claude-3-sonnet-20240229"
)

type AIConfig struct {
	Service       string
	APIKey        string
	AnthropicKey  string
	Model         string
}

func GetAIConfig() AIConfig {
	service := os.Getenv(EnvAIService)
	if service == "" {
		service = DefaultAIService
	}
	
	if service == AIServiceNone {
		return AIConfig{
			Service: AIServiceNone,
		}
	}
	
	apiKey := os.Getenv(EnvAIAPIKey)
	anthropicKey := os.Getenv(EnvAIAnthropicAPIKey)
	
	if service == AIServiceClaude && anthropicKey == "" {
		anthropicKey = apiKey
	}
	
	model := os.Getenv(EnvAIModel)
	if model == "" {
		switch service {
		case AIServiceClaude:
			model = DefaultClaudeModel
		default:
			model = DefaultOpenAIModel
		}
	}
	
	return AIConfig{
		Service:      service,
		APIKey:       apiKey,
		AnthropicKey: anthropicKey,
		Model:        model,
	}
} 