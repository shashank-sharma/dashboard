package connectors

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/shashank-sharma/backend/internal/services/workflow/types"
	"github.com/shashank-sharma/backend/internal/store"
)

// CSVConnector is a connector for reading and writing CSV files
type CSVConnector struct {
	types.BaseConnector
}

// NewCSVSourceConnector creates a new CSV source connector
func NewCSVSourceConnector() types.Connector {
	configSchema := map[string]interface{}{
		"file_path": map[string]interface{}{
			"type":        "string",
			"title":       "File Path",
			"description": "Path to the CSV file (uploads/{filename}.csv for uploaded files)",
			"required":    true,
		},
		"has_header": map[string]interface{}{
			"type":        "boolean",
			"title":       "Has Header",
			"description": "Whether the CSV file has a header row",
			"default":     true,
			"required":    false,
		},
		"delimiter": map[string]interface{}{
			"type":        "string",
			"title":       "Delimiter",
			"description": "Field delimiter (comma, semicolon, tab, etc.)",
			"default":     ",",
			"required":    false,
		},
		"comment": map[string]interface{}{
			"type":        "string",
			"title":       "Comment Character",
			"description": "Character that marks the start of a comment line",
			"required":    false,
		},
	}

	connector := &CSVConnector{
		BaseConnector: types.BaseConnector{
			ConnID:       "csv_source",
			ConnName:     "CSV Source",
			ConnType:     types.SourceConnector,
			ConfigSchema: configSchema,
			Config:       make(map[string]interface{}),
		},
	}

	return connector
}

// NewCSVDestinationConnector creates a new CSV destination connector
func NewCSVDestinationConnector() types.Connector {
	configSchema := map[string]interface{}{
		"file_path": map[string]interface{}{
			"type":        "string",
			"title":       "File Path",
			"description": "Path to the CSV file (uploads/{filename}.csv for uploaded files)",
			"required":    true,
		},
		"delimiter": map[string]interface{}{
			"type":        "string",
			"title":       "Delimiter",
			"description": "Field delimiter (comma, semicolon, tab, etc.)",
			"default":     ",",
			"required":    false,
		},
		"include_header": map[string]interface{}{
			"type":        "boolean",
			"title":       "Include Header",
			"description": "Whether to include a header row in the CSV file",
			"default":     true,
			"required":    false,
		},
	}

	connector := &CSVConnector{
		BaseConnector: types.BaseConnector{
			ConnID:       "csv_destination",
			ConnName:     "CSV Destination",
			ConnType:     types.DestinationConnector,
			ConfigSchema: configSchema,
			Config:       make(map[string]interface{}),
		},
	}

	return connector
}

// Execute runs the CSV connector operation and returns the result
func (c *CSVConnector) Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	config := c.Config
	
	filePath, ok := config["file_path"].(string)
	if !ok {
		return nil, fmt.Errorf("file path is required")
	}
	
	if c.Type() == types.SourceConnector {
		return c.readCSV(filePath, config)
	} else if c.Type() == types.DestinationConnector {
		return c.writeCSV(filePath, config, input)
	}
	
	return nil, fmt.Errorf("unsupported connector type: %s", c.Type())
}

// resolveFilePath resolves the file path based on the storage path
func (c *CSVConnector) resolveFilePath(filePath string) string {
	pb := store.GetDao()
	
	// If it's an uploaded file path that starts with "uploads/"
	if strings.HasPrefix(filePath, "uploads/") {
		return filepath.Join(pb.DataDir(), filePath)
	}
	
	// For destination files, store in the workflow results directory
	if c.Type() == types.DestinationConnector {
		resultsDir := filepath.Join(pb.DataDir(), "storage", "workflow_results")
		
		if err := os.MkdirAll(resultsDir, 0755); err != nil {
			return filePath
		}
		
		return filepath.Join(resultsDir, filePath)
	}
	
	return filePath
}

// readCSV reads data from a CSV file
func (c *CSVConnector) readCSV(filePath string, config map[string]interface{}) (map[string]interface{}, error) {
	resolvedPath := c.resolveFilePath(filePath)
	
	file, err := os.Open(resolvedPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()
	
	reader := csv.NewReader(file)
	
	if delimiter, ok := config["delimiter"].(string); ok && len(delimiter) > 0 {
		reader.Comma = rune(delimiter[0])
	}
	
	if comment, ok := config["comment"].(string); ok && len(comment) > 0 {
		reader.Comment = rune(comment[0])
	}
	
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV data: %w", err)
	}
	
	// Convert to map or array based on header setting
	hasHeader, _ := config["has_header"].(bool)
	
	var result []map[string]interface{}
	
	if len(records) == 0 {
		return map[string]interface{}{
			"data": []interface{}{},
		}, nil
	}
	
	if hasHeader && len(records) > 0 {
		headers := records[0]
		
		for i := 1; i < len(records); i++ {
			row := make(map[string]interface{})
			
			for j, value := range records[i] {
				if j < len(headers) {
					row[headers[j]] = value
				} else {
					row[fmt.Sprintf("column_%d", j+1)] = value
				}
			}
			
			result = append(result, row)
		}
	} else {
		for _, record := range records {
			row := make(map[string]interface{})
			
			for j, value := range record {
				row[fmt.Sprintf("column_%d", j+1)] = value
			}
			
			result = append(result, row)
		}
	}
	
	return map[string]interface{}{
		"data": result,
		"file_path": resolvedPath,
		"record_count": len(result),
	}, nil
}

// writeCSV writes data to a CSV file
func (c *CSVConnector) writeCSV(filePath string, config map[string]interface{}, input map[string]interface{}) (map[string]interface{}, error) {
	resolvedPath := c.resolveFilePath(filePath)
	
	inputData, ok := input["data"]
	if !ok {
		return nil, fmt.Errorf("no data provided for CSV output")
	}
	
	rows, headers, err := convertToCSVRecords(inputData)
	if err != nil {
		return nil, err
	}
	
	append, _ := config["append"].(bool)
	
	var fileMode int
	if append {
		fileMode = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	} else {
		fileMode = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	}
	
	err = os.MkdirAll(filepath.Dir(resolvedPath), 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}
	
	file, err := os.OpenFile(resolvedPath, fileMode, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file for writing: %w", err)
	}
	defer file.Close()
	
	writer := csv.NewWriter(file)
	defer writer.Flush()
	
	if delimiter, ok := config["delimiter"].(string); ok && len(delimiter) > 0 {
		writer.Comma = rune(delimiter[0])
	}
	
	includeHeader, _ := config["include_header"].(bool)
	
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}
	
	if includeHeader && (!append || fileInfo.Size() == 0) {
		if err := writer.Write(headers); err != nil {
			return nil, fmt.Errorf("failed to write CSV header: %w", err)
		}
	}
	
	for _, row := range rows {
		if err := writer.Write(row); err != nil {
			return nil, fmt.Errorf("failed to write CSV row: %w", err)
		}
	}
	
	return map[string]interface{}{
		"file_path":    resolvedPath,
		"record_count": len(rows),
		"success":      true,
	}, nil
}

// convertToCSVRecords converts input data to CSV records
func convertToCSVRecords(input interface{}) ([][]string, []string, error) {
	var rows [][]string
	var headers []string
	
	inputArray, ok := input.([]interface{})
	if ok {
		if len(inputArray) == 0 {
			return [][]string{}, []string{}, nil
		}
		
		headerMap := make(map[string]bool)
		
		for _, item := range inputArray {
			if objItem, ok := item.(map[string]interface{}); ok {
				for key := range objItem {
					headerMap[key] = true
				}
			}
		}
		
		for header := range headerMap {
			headers = append(headers, header)
		}
		
		for _, item := range inputArray {
			if objItem, ok := item.(map[string]interface{}); ok {
				row := make([]string, len(headers))
				
				for i, header := range headers {
					if val, exists := objItem[header]; exists {
						row[i] = fmt.Sprintf("%v", val)
					} else {
						row[i] = ""
					}
				}
				
				rows = append(rows, row)
			}
		}
	} else if objInput, ok := input.(map[string]interface{}); ok {
		if data, ok := objInput["data"].([]interface{}); ok {
			return convertToCSVRecords(data)
		}
		
		for key := range objInput {
			headers = append(headers, key)
		}
		
		row := make([]string, len(headers))
		for i, header := range headers {
			if val, exists := objInput[header]; exists {
				row[i] = fmt.Sprintf("%v", val)
			} else {
				row[i] = ""
			}
		}
		
		rows = append(rows, row)
	} else {
		return nil, nil, fmt.Errorf("input data must be an array of objects or a single object")
	}
	
	return rows, headers, nil
} 