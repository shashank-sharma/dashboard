package connectors

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/shashank-sharma/backend/internal/logger"
	"github.com/shashank-sharma/backend/internal/services/workflow/types"
)

// HTTPConnector is a connector for making HTTP requests
type HTTPConnector struct {
	types.BaseConnector
	client *http.Client
}

// NewHTTPSourceConnector creates a new HTTP source connector
func NewHTTPSourceConnector() types.Connector {
	configSchema := map[string]interface{}{
		"url": map[string]interface{}{
			"type":        "string",
			"title":       "URL",
			"description": "URL to make the request to",
			"required":    true,
		},
		"method": map[string]interface{}{
			"type":        "string",
			"title":       "Method",
			"description": "HTTP method (GET, POST, PUT, DELETE)",
			"enum":        []string{"GET", "POST", "PUT", "DELETE"},
			"default":     "GET",
			"required":    true,
		},
		"headers": map[string]interface{}{
			"type":        "object",
			"title":       "Headers",
			"description": "HTTP headers to include in the request",
			"required":    false,
		},
		"body": map[string]interface{}{
			"type":        "string",
			"title":       "Body",
			"description": "Request body (for POST/PUT requests)",
			"required":    false,
		},
		"params": map[string]interface{}{
			"type":        "object",
			"title":       "Query Parameters",
			"description": "URL query parameters",
			"required":    false,
		},
		"timeout": map[string]interface{}{
			"type":        "integer",
			"title":       "Timeout",
			"description": "Request timeout in seconds",
			"default":     30,
			"required":    false,
		},
		"parse_json": map[string]interface{}{
			"type":        "boolean",
			"title":       "Parse JSON",
			"description": "Parse the response as JSON",
			"default":     true,
			"required":    false,
		},
	}

	connector := &HTTPConnector{
		BaseConnector: types.BaseConnector{
			ConnID:       "http",
			ConnName:     "HTTP Request",
			ConnType:     types.SourceConnector,
			ConfigSchema: configSchema,
			Config:       make(map[string]interface{}),
		},
		client: &http.Client{},
	}

	return connector
}

// NewHTTPDestinationConnector creates a new HTTP destination connector
func NewHTTPDestinationConnector() types.Connector {
	// Reuse the source connector's config schema
	connector := NewHTTPSourceConnector().(*HTTPConnector)
	
	connector.ConnID = "http_destination"
	connector.ConnName = "HTTP Destination"
	connector.ConnType = types.DestinationConnector
	
	return connector
}

// Configure sets up the connector with the provided configuration
func (c *HTTPConnector) Configure(config map[string]interface{}) error {
	if err := c.BaseConnector.Configure(config); err != nil {
		return err
	}

	// Set up HTTP client with timeout
	timeout := 30 * time.Second
	if timeoutVal, ok := config["timeout"]; ok {
		if t, ok := timeoutVal.(float64); ok {
			timeout = time.Duration(t) * time.Second
		}
	}
	c.client.Timeout = timeout

	return nil
}

// Execute runs the HTTP connector
func (c *HTTPConnector) Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	// Get configuration
	config := c.Config
	
	// Get the request URL
	requestURL, ok := config["url"].(string)
	if !ok || requestURL == "" {
		return nil, fmt.Errorf("URL is required")
	}
	
	// Get the request method
	method, _ := config["method"].(string)
	if method == "" {
		method = "GET"
	}
	
	// Add query parameters if any
	if params, ok := config["params"].(map[string]interface{}); ok && len(params) > 0 {
		parsedURL, err := url.Parse(requestURL)
		if err != nil {
			return nil, fmt.Errorf("invalid URL: %w", err)
		}
		
		query := parsedURL.Query()
		for key, value := range params {
			query.Add(key, fmt.Sprintf("%v", value))
		}
		
		parsedURL.RawQuery = query.Encode()
		requestURL = parsedURL.String()
	}
	
	// Prepare request body if needed
	var bodyReader io.Reader
	if body, ok := config["body"].(string); ok && body != "" {
		bodyReader = strings.NewReader(body)
	}
	
	// Create the request
	req, err := http.NewRequestWithContext(ctx, method, requestURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	// Add headers
	if headers, ok := config["headers"].(map[string]interface{}); ok {
		for key, value := range headers {
			req.Header.Add(key, fmt.Sprintf("%v", value))
		}
	}
	
	// Set timeout if specified
	timeout := 30
	if timeoutVal, ok := config["timeout"].(float64); ok {
		timeout = int(timeoutVal)
	}
	c.client.Timeout = time.Duration(timeout) * time.Second
	
	// Execute the request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	
	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	
	// Parse response as JSON if requested
	parseJSON := true
	if parseJSONVal, ok := config["parse_json"].(bool); ok {
		parseJSON = parseJSONVal
	}
	
	var responseData interface{} = string(responseBody)
	if parseJSON && len(responseBody) > 0 {
		var jsonData interface{}
		if err := json.Unmarshal(responseBody, &jsonData); err == nil {
			responseData = jsonData
		} else {
			logger.Info.Printf("Failed to parse response as JSON: %v", err)
		}
	}
	
	// Prepare the result
	result := map[string]interface{}{
		"status_code": resp.StatusCode,
		"status_text": resp.Status,
		"headers":     resp.Header,
		"body":        responseData,
	}
	
	return result, nil
} 