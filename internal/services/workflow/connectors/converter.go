package connectors

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/shashank-sharma/backend/internal/logger"
	"github.com/shashank-sharma/backend/internal/services/workflow/types"
)

// PBToCsvConverter is a utility connector that converts PocketBase data to CSV
type PBToCsvConverter struct {
	types.BaseConnector
}

// NewPBToCsvConverter creates a new PocketBase to CSV converter connector
func NewPBToCsvConverter() types.Connector {
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"inputFormat": map[string]interface{}{
				"type":        "string",
				"title":       "Input Format",
				"description": "Format of the input data (json)",
				"default":     "json",
				"enum":        []string{"json"},
			},
			"outputPath": map[string]interface{}{
				"type":        "string",
				"title":       "Output Path",
				"description": "Path where the CSV file will be saved (relative to storage folder)",
				"required":    true,
			},
			"includeHeader": map[string]interface{}{
				"type":        "boolean",
				"title":       "Include Header",
				"description": "Whether to include a header row in the CSV file",
				"default":     true,
			},
		},
	}

	connector := &PBToCsvConverter{
		BaseConnector: types.BaseConnector{
			ConnID:       "pb_to_csv_converter",
			ConnName:     "PocketBase to CSV Converter",
			ConnType:     types.ProcessorConnector,
			ConfigSchema: schema,
			Config:       make(map[string]interface{}),
		},
	}
	return connector
}

// Execute converts the input data to a CSV file
func (c *PBToCsvConverter) Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	outputPath := c.Config["outputPath"].(string)
	includeHeader := true
	if val, exists := c.Config["includeHeader"]; exists {
		if boolVal, ok := val.(bool); ok {
			includeHeader = boolVal
		}
	}
	
	if input == nil || len(input) == 0 {
		return nil, fmt.Errorf("no data provided for CSV conversion")
	}
	
	var records []map[string]interface{}
	
	// Try to extract records from different input formats
	if data, exists := input["data"]; exists {
		if recordsArray, ok := data.([]map[string]interface{}); ok {
			records = recordsArray
		} else if recordsMap, ok := data.(map[string]interface{}); ok {
			for _, v := range recordsMap {
				if array, ok := v.([]map[string]interface{}); ok {
					records = array
					break
				}
			}
		}
	}
	
	if records == nil {
		if array, ok := input["records"].([]map[string]interface{}); ok {
			records = array
		}
	}
	
	// If still no records, try to extract from JSON string
	if records == nil {
		if jsonStr, ok := input["json"].(string); ok {
			var data map[string]interface{}
			if err := json.Unmarshal([]byte(jsonStr), &data); err == nil {
				// Try to find an array in the unmarshaled data
				for _, v := range data {
					if array, ok := v.([]map[string]interface{}); ok {
						records = array
						break
					}
				}
			}
		}
	}
	
	// If still no records, return error
	if records == nil || len(records) == 0 {
		return nil, fmt.Errorf("no valid records found in input data")
	}
	
	// Ensure output directory exists
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}
	
	// Convert to CSV
	if err := convertToCsv(records, outputPath, includeHeader); err != nil {
		return nil, fmt.Errorf("failed to convert to CSV: %w", err)
	}
	
	return map[string]interface{}{
		"status":      "success",
		"file_path":   outputPath,
		"record_count": len(records),
	}, nil
}

// convertToCsv converts a slice of records to a CSV file
func convertToCsv(records []map[string]interface{}, outputPath string, includeHeader bool) error {
	if len(records) == 0 {
		return fmt.Errorf("no records to convert")
	}
	
	// Create CSV file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()
	
	writer := csv.NewWriter(file)
	defer writer.Flush()
	
	// Extract headers from the first record
	var headers []string
	for key := range records[0] {
		headers = append(headers, key)
	}
	
	// Write header if required
	if includeHeader {
		if err := writer.Write(headers); err != nil {
			return fmt.Errorf("failed to write CSV header: %w", err)
		}
	}
	
	// Write records
	for _, record := range records {
		row := make([]string, len(headers))
		for i, header := range headers {
			// Convert value to string
			value := record[header]
			if value == nil {
				row[i] = ""
			} else {
				row[i] = fmt.Sprintf("%v", value)
			}
		}
		
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}
	
	logger.Info.Printf("Successfully converted %d records to CSV at %s", len(records), outputPath)
	return nil
} 