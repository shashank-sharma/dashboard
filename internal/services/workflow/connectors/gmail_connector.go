package connectors

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/shashank-sharma/backend/internal/services/workflow/types"
	"github.com/shashank-sharma/backend/internal/store"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// GmailConnector is a connector for reading emails from Gmail
// TODO: Not required, used for testing
type GmailConnector struct {
	types.BaseConnector
	client *gmail.Service
}

// NewGmailSourceConnector creates a new Gmail source connector
func NewGmailSourceConnector() types.Connector {
	configSchema := map[string]interface{}{
		"token_id": map[string]interface{}{
			"type":        "string",
			"title":       "Token ID",
			"description": "ID of the saved Gmail token in the tokens collection",
			"required":    true,
		},
		"query": map[string]interface{}{
			"type":        "string",
			"title":       "Query",
			"description": "Gmail search query (e.g. 'from:someone@example.com is:unread')",
			"default":     "is:unread",
			"required":    false,
		},
		"max_results": map[string]interface{}{
			"type":        "number",
			"title":       "Max Results",
			"description": "Maximum number of emails to fetch",
			"default":     10,
			"required":    false,
		},
		"include_content": map[string]interface{}{
			"type":        "boolean",
			"title":       "Include Content",
			"description": "Whether to include the email content in the results",
			"default":     true,
			"required":    false,
		},
	}

	connector := &GmailConnector{
		BaseConnector: types.BaseConnector{
			ConnID:       "gmail_source",
			ConnName:     "Gmail Source",
			ConnType:     types.SourceConnector,
			ConfigSchema: configSchema,
			Config:       make(map[string]interface{}),
		},
	}

	return connector
}

// Configure configures the connector and initializes the Gmail client
func (c *GmailConnector) Configure(config map[string]interface{}) error {
	// Call base Configure method
	if err := c.BaseConnector.Configure(config); err != nil {
		return err
	}

	// Initialize Gmail client if token ID is provided
	if tokenID, ok := config["token_id"].(string); ok && tokenID != "" {
		if err := c.initializeGmailClient(tokenID); err != nil {
			return fmt.Errorf("failed to initialize Gmail client: %w", err)
		}
	}

	return nil
}

// initializeGmailClient initializes the Gmail API client using the token from the database
func (c *GmailConnector) initializeGmailClient(tokenID string) error {
	// Get token from the database using the store package
	record, err := store.GetDao().FindRecordById("tokens", tokenID)
	if err != nil {
		return fmt.Errorf("failed to find token with ID %s: %w", tokenID, err)
	}

	// Extract token data
	accessToken := record.GetString("access_token")
	refreshToken := record.GetString("refresh_token")
	expiry := record.GetDateTime("expiry")
	tokenType := record.GetString("token_type")

	// Check if the token is valid
	if accessToken == "" || refreshToken == "" {
		return fmt.Errorf("invalid token data")
	}

	// Create OAuth2 token
	token := &oauth2.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    tokenType,
		Expiry:       expiry.Time(),
	}

	// Configure OAuth2 config
	config := &oauth2.Config{
		ClientID:     "YOUR_CLIENT_ID",
		ClientSecret: "YOUR_CLIENT_SECRET",
		Endpoint:     google.Endpoint,
		Scopes:       []string{gmail.GmailReadonlyScope},
	}

	// Create OAuth2 client
	client := config.Client(context.Background(), token)

	// Create Gmail service
	gmailService, err := gmail.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("failed to create Gmail service: %w", err)
	}

	c.client = gmailService
	return nil
}

// Execute runs the Gmail connector operation and returns the emails
func (c *GmailConnector) Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	config := c.Config

	// Initialize Gmail client if not done yet
	if c.client == nil {
		tokenID, ok := config["token_id"].(string)
		if !ok {
			return nil, fmt.Errorf("token ID is required")
		}

		if err := c.initializeGmailClient(tokenID); err != nil {
			return nil, err
		}
	}

	// Get parameters from config
	query, _ := config["query"].(string)
	maxResults := 10
	if maxVal, ok := config["max_results"].(float64); ok {
		maxResults = int(maxVal)
	}
	includeBody, _ := config["include_content"].(bool)

	// Get label IDs if specified
	var labelIDs []string
	if labelIDsConfig, ok := config["label_ids"].([]interface{}); ok {
		for _, labelIDInterface := range labelIDsConfig {
			if labelID, ok := labelIDInterface.(string); ok {
				labelIDs = append(labelIDs, labelID)
			}
		}
	}

	// Create the list request
	listReq := c.client.Users.Messages.List("me").MaxResults(int64(maxResults))
	
	// Add query if specified
	if query != "" {
		listReq = listReq.Q(query)
	}
	
	// Add label IDs if specified
	if len(labelIDs) > 0 {
		listReq = listReq.LabelIds(labelIDs...)
	}

	// Execute the request
	resp, err := listReq.Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list messages: %w", err)
	}

	// Process the messages
	emails := make([]map[string]interface{}, 0)
	
	for _, message := range resp.Messages {
		// Get the full message
		msg, err := c.client.Users.Messages.Get("me", message.Id).Do()
		if err != nil {
			return nil, fmt.Errorf("failed to get message %s: %w", message.Id, err)
		}

		// Extract headers
		headers := make(map[string]string)
		for _, header := range msg.Payload.Headers {
			headers[strings.ToLower(header.Name)] = header.Value
		}

		// Create email object
		email := map[string]interface{}{
			"id":       msg.Id,
			"thread_id": msg.ThreadId,
			"from":     headers["from"],
			"to":       headers["to"],
			"subject":  headers["subject"],
			"date":     headers["date"],
			"snippet":  msg.Snippet,
			"labels":   msg.LabelIds,
		}

		// Add body if requested
		if includeBody {
			body, err := extractMessageBody(msg)
			if err != nil {
				return nil, fmt.Errorf("failed to extract body from message %s: %w", message.Id, err)
			}
			email["body"] = body
		}

		emails = append(emails, email)
	}

	// Return the results
	return map[string]interface{}{
		"data":         emails,
		"result_count": len(emails),
		"has_more":     resp.NextPageToken != "",
		"next_page_token": resp.NextPageToken,
	}, nil
}

// extractMessageBody extracts the body from a Gmail message
func extractMessageBody(msg *gmail.Message) (string, error) {
	if msg.Payload == nil {
		return "", nil
	}

	// For simple messages
	if msg.Payload.Body != nil && msg.Payload.Body.Data != "" {
		data, err := base64.URLEncoding.DecodeString(msg.Payload.Body.Data)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}

	// For multipart messages
	if len(msg.Payload.Parts) > 0 {
		for _, part := range msg.Payload.Parts {
			if part.MimeType == "text/plain" || part.MimeType == "text/html" {
				if part.Body != nil && part.Body.Data != "" {
					data, err := base64.URLEncoding.DecodeString(part.Body.Data)
					if err != nil {
						return "", err
					}
					return string(data), nil
				}
			}
		}
	}

	return "", nil
} 