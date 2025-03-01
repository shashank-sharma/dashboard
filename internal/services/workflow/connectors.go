package workflow

import (
	"context"
	"fmt"

	"github.com/shashank-sharma/backend/internal/services/workflow/types"
)

// Simple implementations of various connectors

// FileSourceConnector reads data from a file
type FileSourceConnector struct {
	types.BaseConnector
}

// Execute reads data from a file or directory
func (c *FileSourceConnector) Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	// Simple implementation - in real code this would access files
	return map[string]interface{}{
		"data": "File content would be read here",
		"path": c.Config["filePath"],
	}, nil
}

// NewFileSourceConnector creates a new file source connector
func NewFileSourceConnector() types.Connector {
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"filePath": map[string]interface{}{
				"type":        "string",
				"description": "Path to the file or directory",
			},
		},
		"required": []string{"filePath"},
	}

	connector := &FileSourceConnector{
		BaseConnector: types.BaseConnector{
			ConnID:       "file_source",
			ConnName:     "File Source",
			ConnType:     types.SourceConnector,
			ConfigSchema: schema,
			Config:       make(map[string]interface{}),
		},
	}
	return connector
}

// LogDestinationConnector outputs data to logs
type LogDestinationConnector struct {
	types.BaseConnector
}

// Execute outputs data to logs
func (c *LogDestinationConnector) Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	// Simple implementation - in real code this would log to appropriate destination
	message := fmt.Sprintf("Log: %v", input)
	return map[string]interface{}{
		"status":  "success",
		"message": message,
	}, nil
}

// NewLogDestinationConnector creates a new log destination connector
func NewLogDestinationConnector() types.Connector {
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"logLevel": map[string]interface{}{
				"type":    "string",
				"enum":    []string{"debug", "info", "warn", "error"},
				"default": "info",
			},
		},
	}

	connector := &LogDestinationConnector{
		BaseConnector: types.BaseConnector{
			ConnID:       "log_destination",
			ConnName:     "Log Destination",
			ConnType:     types.DestinationConnector,
			ConfigSchema: schema,
			Config:       make(map[string]interface{}),
		},
	}
	return connector
}

// TransformProcessor transforms input data
type TransformProcessor struct {
	types.BaseConnector
}

// Execute transforms input data according to configuration
func (c *TransformProcessor) Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	// Simple implementation - in real code this would apply transformations
	return map[string]interface{}{
		"transformed": true,
		"original":    input,
	}, nil
}

// NewTransformProcessor creates a new transform processor connector
func NewTransformProcessor() types.Connector {
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"transformations": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"field":  map[string]interface{}{"type": "string"},
						"action": map[string]interface{}{"type": "string", "enum": []string{"rename", "delete", "modify"}},
						"value":  map[string]interface{}{"type": "string"},
					},
				},
			},
		},
	}

	connector := &TransformProcessor{
		BaseConnector: types.BaseConnector{
			ConnID:       "transform_processor",
			ConnName:     "Transform Processor",
			ConnType:     types.ProcessorConnector,
			ConfigSchema: schema,
			Config:       make(map[string]interface{}),
		},
	}
	return connector
}