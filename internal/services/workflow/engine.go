package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	pbTypes "github.com/pocketbase/pocketbase/tools/types"
	"github.com/shashank-sharma/backend/internal/logger"
	wfModels "github.com/shashank-sharma/backend/internal/models"
	"github.com/shashank-sharma/backend/internal/query"
	"github.com/shashank-sharma/backend/internal/services/workflow/connectors"
	"github.com/shashank-sharma/backend/internal/services/workflow/types"
	"github.com/shashank-sharma/backend/internal/store"
)

// Node represents a node in a workflow graph
type Node struct {
	ID       string                 `json:"id"`       // Unique identifier for the node
	Name     string                 `json:"name"`     // Display name of the node
	Type     string                 `json:"type"`     // Category of the node: "source", "processor", or "destination"
	NodeType string                 `json:"node_type"` // Specific connector type (e.g., "pocketbase_source", "pb_to_csv_converter")
	Config   map[string]interface{} `json:"config"`   // Configuration parameters for the node
	Inputs   []string               `json:"inputs,omitempty"`  // IDs of nodes that feed into this node
	Outputs  []string               `json:"outputs,omitempty"` // IDs of nodes that this node feeds into
	Position struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	} `json:"position"` // Visual position of the node in the workflow editor
}

// Edge represents a connection between nodes
type Edge struct {
	ID     string `json:"id"`     // Unique identifier for the edge
	Source string `json:"source"` // ID of the source node
	Target string `json:"target"` // ID of the target node
}

// WorkflowData represents the structure of a workflow
type WorkflowData struct {
	Nodes []*Node `json:"nodes"` // List of nodes in the workflow
	Edges []*Edge `json:"edges"` // List of connections between nodes
}

// Graph represents the workflow execution graph
type Graph struct {
	Nodes map[string]*Node // Map of node ID to node
	Edges []*Edge          // List of connections between nodes
}

// WorkflowEngine handles the execution of workflows
type WorkflowEngine struct {
	dao      *pocketbase.PocketBase
	registry types.ConnectorRegistry
}

// NewWorkflowEngine creates a new workflow engine
func NewWorkflowEngine(pb *pocketbase.PocketBase) *WorkflowEngine {
	registry := NewConnectorRegistryImpl()
	
	// Register all connectors using the centralized function
	RegisterConnectors(registry)
	
	// Log all available connectors for debugging
	availableConnectors := registry.GetAvailableConnectors()
	logger.Info.Printf("Available workflow connectors: %v", availableConnectors)
	
	return &WorkflowEngine{
		dao:      pb,
		registry: registry,
	}
}

// ExecuteWorkflow executes a workflow by its ID
func (e *WorkflowEngine) ExecuteWorkflow(ctx context.Context, workflowID string) (*wfModels.WorkflowExecution, error) {
	workflowExecution := &wfModels.WorkflowExecution{
		WorkflowID: workflowID,
		Status: "running",
		StartTime: pbTypes.NowDateTime(),
		Logs: "[]",
		Results: "{}",
	}

	err := query.UpsertRecord[*wfModels.WorkflowExecution](workflowExecution, map[string]interface{}{
		"workflow_id": workflowID,
		"status": "running",
		"start_time": pbTypes.NowDateTime(),
		"logs": "[]",
		"results": "{}",
	})

	if err != nil {
		logger.Error.Println("Failed updating record", err)
	}

	go func() {
		e.runWorkflow(ctx, workflowID, workflowExecution.Id)
	}()

	return workflowExecution, nil
}

// GetExecutionStatus retrieves the status of a workflow execution
func (e *WorkflowEngine) GetExecutionStatus(executionID string) (map[string]interface{}, error) {
	execution, err := query.FindById[*wfModels.WorkflowExecution](executionID)
	if err != nil {
		return nil, fmt.Errorf("execution not found: %w", err)
	}

	status := map[string]interface{}{
		"id":             execution.Id,
		"workflow_id":    execution.WorkflowID,
		"status":         execution.Status,
		"start_time":     execution.StartTime,
		"end_time":       execution.EndTime,
		"duration":       execution.Duration,
		"error_message":  execution.ErrorMessage,
		"result_file_ids": execution.ResultFileIDs,
	}

	// Parse logs if not empty
	if execution.Logs != "" && execution.Logs != "[]" {
		var logs []interface{}
		if err := json.Unmarshal([]byte(execution.Logs), &logs); err == nil {
			status["logs"] = logs
		} else {
			status["logs"] = []interface{}{}
		}
	} else {
		status["logs"] = []interface{}{}
	}

	// Parse results if not empty
	if execution.Results != "" && execution.Results != "{}" {
		var results map[string]interface{}
		if err := json.Unmarshal([]byte(execution.Results), &results); err == nil {
			status["results"] = results
		} else {
			status["results"] = map[string]interface{}{}
		}
	} else {
		status["results"] = map[string]interface{}{}
	}

	return status, nil
}

// runWorkflow executes the workflow and updates its status
func (e *WorkflowEngine) runWorkflow(ctx context.Context, workflowID string, executionID string) {
	startTime := time.Now()
	logs := make([]map[string]interface{}, 0)
	
	// Add log entry
	logs = append(logs, map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"level":     "info",
		"message":   fmt.Sprintf("Starting workflow execution %s", executionID),
	})

	logger.LogInfo("Starting workflow execution", "executionID", executionID)
	e.updateExecutionLogs(executionID, logs)

	// Load workflow
	workflow, err := e.loadWorkflow(workflowID)
	if err != nil {
		e.updateExecutionStatus(executionID, "failed", err.Error(), startTime, logs, nil)
		return
	}

	// Log workflow loaded
	logs = append(logs, map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"level":     "info",
		"message":   fmt.Sprintf("Loaded workflow: %s", workflow.GetString("name")),
	})
	e.updateExecutionLogs(executionID, logs)

	// Load workflow nodes and connections
	logger.LogInfo("Loading workflow nodes and connections", "workflowID", workflowID)
	nodes, err := e.getWorkflowNodes(workflowID)
	if err != nil {
		logs = append(logs, map[string]interface{}{
			"timestamp": time.Now().Format(time.RFC3339),
			"level":     "error",
			"message":   fmt.Sprintf("Failed to load workflow nodes: %v", err),
		})
		e.updateExecutionStatus(executionID, "failed", err.Error(), startTime, logs, nil)
		return
	}

	connections, err := e.loadWorkflowConnections(workflowID)
	if err != nil {
		logs = append(logs, map[string]interface{}{
			"timestamp": time.Now().Format(time.RFC3339),
			"level":     "error",
			"message":   fmt.Sprintf("Failed to load workflow connections: %v", err),
		})
		e.updateExecutionStatus(executionID, "failed", err.Error(), startTime, logs, nil)
		return
	}

	// Log nodes and connections loaded
	logs = append(logs, map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"level":     "info",
		"message":   fmt.Sprintf("Loaded %d nodes and %d connections", len(nodes), len(connections)),
	})
	e.updateExecutionLogs(executionID, logs)

	// Build execution graph
	graph, err := e.buildGraph(nodes, connections)
	if err != nil {
		logs = append(logs, map[string]interface{}{
			"timestamp": time.Now().Format(time.RFC3339),
			"level":     "error",
			"message":   fmt.Sprintf("Failed to build execution graph: %v", err),
		})
		e.updateExecutionStatus(executionID, "failed", err.Error(), startTime, logs, nil)
		return
	}

	// Log graph built
	logs = append(logs, map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"level":     "info",
		"message":   "Built execution graph",
	})
	e.updateExecutionLogs(executionID, logs)

	// Execute graph
	results, nodeResults, err := e.executeGraph(ctx, graph, logs)
	if err != nil {
		logs = append(logs, map[string]interface{}{
			"timestamp": time.Now().Format(time.RFC3339),
			"level":     "error",
			"message":   fmt.Sprintf("Error executing workflow: %v", err),
		})
		e.updateExecutionStatus(executionID, "failed", err.Error(), startTime, logs, nodeResults)
		return
	}

	// Update execution as completed
	logs = append(logs, map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"level":     "info",
		"message":   "Workflow execution completed successfully",
	})
	e.updateExecutionStatus(executionID, "completed", "", startTime, logs, results)
}

// loadWorkflow loads a workflow from the database
func (e *WorkflowEngine) loadWorkflow(workflowID string) (*core.Record, error) {
	workflow, err := query.FindById[*wfModels.Workflow](workflowID)
	if err != nil {
		return nil, fmt.Errorf("workflow not found: %w", err)
	}

	if !workflow.Active {
		return nil, fmt.Errorf("workflow is not active")
	}

	// Convert to core.Record for backward compatibility
	// This is a temporary solution until the engine is fully refactored to use models directly
	pb := store.GetDao()
	record, err := pb.FindRecordById("workflows", workflowID)
	if err != nil {
		return nil, fmt.Errorf("record not found: %w", err)
	}

	return record, nil
}

// getWorkflowNodes loads all nodes for a workflow
func (e *WorkflowEngine) getWorkflowNodes(workflowID string) ([]*core.Record, error) {
	nodes, err := query.FindAllByFilter[*wfModels.WorkflowNode](map[string]interface{}{
		"workflow_id": workflowID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query workflow nodes: %w", err)
	}

	// Convert to core.Record for backward compatibility
	// This is a temporary solution until the engine is fully refactored to use models directly
	records := make([]*core.Record, len(nodes))
	for i, node := range nodes {
		pb := store.GetDao()
		coreRecord, err := pb.FindRecordById("workflow_nodes", node.Id)
		if err != nil {
			return nil, fmt.Errorf("failed to convert to core.Record: %w", err)
		}
		records[i] = coreRecord
	}

	return records, nil
}

// loadWorkflowConnections loads all connections for a workflow
func (e *WorkflowEngine) loadWorkflowConnections(workflowID string) ([]*core.Record, error) {
	connections, err := query.FindAllByFilter[*wfModels.WorkflowConnection](map[string]interface{}{
		"workflow_id": workflowID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query workflow connections: %w", err)
	}

	// Convert to core.Record for backward compatibility
	// This is a temporary solution until the engine is fully refactored to use models directly
	records := make([]*core.Record, len(connections))
	for i, conn := range connections {
		pb := store.GetDao()
		coreRecord, err := pb.FindRecordById("workflow_connections", conn.Id)
		if err != nil {
			return nil, fmt.Errorf("failed to convert to core.Record: %w", err)
		}
		records[i] = coreRecord
	}

	return records, nil
}

// buildGraph builds the execution graph from nodes and connections
func (e *WorkflowEngine) buildGraph(nodes []*core.Record, connections []*core.Record) (*Graph, error) {
	graph := &Graph{
		Nodes: make(map[string]*Node),
		Edges: make([]*Edge, 0, len(connections)),
	}

	// Create nodes
	for _, node := range nodes {
		configJSON := node.GetString("config")
		var config map[string]interface{}

		if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
			return nil, fmt.Errorf("invalid config for node %s: %w", node.Id, err)
		}

		graph.Nodes[node.Id] = &Node{
			ID:       node.Id,
			Name:     node.GetString("name"),
			Type:     node.GetString("type"),
			NodeType: node.GetString("node_type"),
			Config:   config,
			Inputs:   []string{},
			Outputs:  []string{},
			Position: struct {
				X float64 `json:"x"`
				Y float64 `json:"y"`
			}{
				X: node.GetFloat("position_x"),
				Y: node.GetFloat("position_y"),
			},
		}
	}

	// Create edges
	for _, conn := range connections {
		sourceID := conn.GetString("source_id")
		targetID := conn.GetString("target_id")

		if _, exists := graph.Nodes[sourceID]; !exists {
			return nil, fmt.Errorf("source node %s not found", sourceID)
		}

		if _, exists := graph.Nodes[targetID]; !exists {
			return nil, fmt.Errorf("target node %s not found", targetID)
		}

		edge := &Edge{
			ID:     conn.Id,
			Source: sourceID,
			Target: targetID,
		}

		graph.Edges = append(graph.Edges, edge)

		// Update node inputs and outputs
		graph.Nodes[sourceID].Outputs = append(graph.Nodes[sourceID].Outputs, targetID)
		graph.Nodes[targetID].Inputs = append(graph.Nodes[targetID].Inputs, sourceID)
	}

	// Validate graph (ensure at least one source and one destination)
	hasSource := false
	hasDestination := false

	for _, coreNode := range nodes {
		nodeType := coreNode.GetString("type")
		if nodeType == "source" {
			hasSource = true
		} else if nodeType == "destination" {
			hasDestination = true
		}
	}

	if !hasSource {
		return nil, fmt.Errorf("workflow must have at least one source node")
	}

	if !hasDestination {
		return nil, fmt.Errorf("workflow must have at least one destination node")
	}

	return graph, nil
}

// executeGraph executes the workflow graph
func (e *WorkflowEngine) executeGraph(ctx context.Context, graph *Graph, logs []map[string]interface{}) (map[string]interface{}, map[string]interface{}, error) {
	nodeResults := make(map[string]interface{})
	
	// Find source nodes (nodes with no inputs)
	sourceNodes := make([]*Node, 0)
	for _, node := range graph.Nodes {
		logger.LogInfo("Node found in graph", 
			"id", node.ID, 
			"name", node.Name, 
			"type", node.Type, 
			"connector", node.NodeType)
			
		if len(node.Inputs) == 0 && node.Type == "source" {
			sourceNodes = append(sourceNodes, node)
		}
	}

	if len(sourceNodes) == 0 {
		return nil, nodeResults, fmt.Errorf("no source nodes found in workflow")
	}

	logger.LogInfo("Found source nodes", "count", len(sourceNodes))
	
	// Process nodes in topological order
	visited := make(map[string]bool)
	finalResults := make(map[string]interface{})

	// Process each source node
	for _, sourceNode := range sourceNodes {
		logs = append(logs, map[string]interface{}{
			"timestamp": time.Now().Format(time.RFC3339),
			"level":     "info",
			"message":   fmt.Sprintf("Executing source node: %s (Type: %s, Connector: %s)", sourceNode.ID, sourceNode.Type, sourceNode.NodeType),
			"node_id":   sourceNode.ID,
			"node_type": sourceNode.Type,
			"connector": sourceNode.NodeType,
		})

		logger.LogInfo("Executing source node", 
			"id", sourceNode.ID, 
			"name", sourceNode.Name, 
			"connector", sourceNode.NodeType)

		result, err := e.executeNode(ctx, sourceNode, nil)
		if err != nil {
			logger.LogError("Failed to execute source node", 
				"id", sourceNode.ID, 
				"error", err.Error())
			return nil, nodeResults, fmt.Errorf("failed to execute source node %s: %w", sourceNode.ID, err)
		}

		nodeResults[sourceNode.ID] = result

		logs = append(logs, map[string]interface{}{
			"timestamp": time.Now().Format(time.RFC3339),
			"level":     "info",
			"message":   fmt.Sprintf("Source node %s executed successfully", sourceNode.ID),
			"node_id":   sourceNode.ID,
			"status":    "success",
		})

		logger.LogInfo("Source node executed successfully", 
			"id", sourceNode.ID)

		if err := e.processOutgoingEdges(ctx, graph, sourceNode.ID, result, visited, nodeResults, logs); err != nil {
			return nil, nodeResults, err
		}
	}

	// Collect results from destination nodes
	destinationCount := 0
	for id, node := range graph.Nodes {
		if node.Type == "destination" && visited[id] {
			if result, exists := nodeResults[id]; exists {
				finalResults[id] = result
				destinationCount++
			}
		}
	}
	
	logger.LogInfo("Workflow execution completed", 
		"source_nodes", len(sourceNodes),
		"destination_nodes", destinationCount,
		"total_nodes_executed", len(visited))

	return finalResults, nodeResults, nil
}

// processOutgoingEdges processes all outgoing edges from a node
func (e *WorkflowEngine) processOutgoingEdges(
	ctx context.Context,
	graph *Graph,
	nodeID string,
	nodeResult interface{},
	visited map[string]bool,
	nodeResults map[string]interface{},
	logs []map[string]interface{},
) error {
	visited[nodeID] = true
	node := graph.Nodes[nodeID]
	
	// If node has no outputs, we're done
	if len(node.Outputs) == 0 {
		logger.LogInfo("Node has no outputs", "id", nodeID, "type", node.Type)
		return nil
	}
	
	logger.LogInfo("Processing outgoing edges", 
		"node_id", nodeID, 
		"output_count", len(node.Outputs))

	// Process each outgoing edge
	for _, targetID := range node.Outputs {
		targetNode := graph.Nodes[targetID]
		
		logger.LogInfo("Checking target node", 
			"source_id", nodeID, 
			"target_id", targetID,
			"target_type", targetNode.Type,
			"target_connector", targetNode.NodeType)

		// Check if all inputs for the target node have been processed
		allInputsProcessed := true
		for _, inputID := range targetNode.Inputs {
			if !visited[inputID] {
				logger.LogInfo("Target node waiting for input", 
					"target_id", targetID, 
					"input_id", inputID)
				allInputsProcessed = false
				break
			}
		}

		if !allInputsProcessed {
			logger.LogInfo("Skipping target node (not all inputs processed)", 
				"target_id", targetID)
			continue
		}

		// Collect inputs from all input nodes
		var input map[string]interface{}
		if lastResult, ok := nodeResult.(map[string]interface{}); ok {
			input = lastResult
		} else {
			// Convert result to map if it's not already
			input = map[string]interface{}{
				"data": nodeResult,
			}
		}

		// Log target node execution
		logs = append(logs, map[string]interface{}{
			"timestamp": time.Now().Format(time.RFC3339),
			"level":     "info",
			"message":   fmt.Sprintf("Executing node: %s (Type: %s, Connector: %s)", targetNode.ID, targetNode.Type, targetNode.NodeType),
			"node_id":   targetNode.ID,
			"node_type": targetNode.Type,
			"connector": targetNode.NodeType,
		})
		
		logger.LogInfo("Executing target node", 
			"id", targetNode.ID, 
			"type", targetNode.Type, 
			"connector", targetNode.NodeType)

		// Execute the target node
		result, err := e.executeNode(ctx, targetNode, input)
		if err != nil {
			logger.LogError("Failed to execute node", 
				"id", targetNode.ID, 
				"error", err.Error())
			return fmt.Errorf("failed to execute node %s: %w", targetNode.ID, err)
		}

		nodeResults[targetNode.ID] = result

		logs = append(logs, map[string]interface{}{
			"timestamp": time.Now().Format(time.RFC3339),
			"level":     "info",
			"message":   fmt.Sprintf("Node %s executed successfully", targetNode.ID),
			"node_id":   targetNode.ID,
			"status":    "success",
		})
		
		logger.LogInfo("Node executed successfully", 
			"id", targetNode.ID)

		if err := e.processOutgoingEdges(ctx, graph, targetNode.ID, result, visited, nodeResults, logs); err != nil {
			return err
		}
	}

	return nil
}

// executeNode executes a single node in the workflow
func (e *WorkflowEngine) executeNode(ctx context.Context, node *Node, input map[string]interface{}) (map[string]interface{}, error) {
	if node.NodeType == "" {
		logger.LogError("Node type is empty", "node_id", node.ID)
		return nil, fmt.Errorf("node type is empty for node %s", node.ID)
	}

	logger.LogInfo("Creating connector", 
		"node_id", node.ID,
		"connector_type", node.NodeType)

	// Create connector instance
	connector, err := e.registry.Create(node.NodeType)
	if err != nil {
		logger.LogError("Failed to create connector", 
			"node_id", node.ID, 
			"connector_type", node.NodeType, 
			"error", err.Error())
		return nil, fmt.Errorf("failed to create connector for node %s: %w", node.ID, err)
	}

	logger.LogInfo("Configuring connector", 
		"node_id", node.ID, 
		"connector_type", node.NodeType)

	// Configure the connector
	if err := connector.Configure(node.Config); err != nil {
		logger.LogError("Failed to configure connector", 
			"node_id", node.ID, 
			"connector_type", node.NodeType, 
			"error", err.Error())
		return nil, fmt.Errorf("failed to configure connector for node %s: %w", node.ID, err)
	}

	// Log input data size for debugging
	var inputSize int
	if input != nil {
		if data, ok := input["data"]; ok {
			// If data is a string, log its length
			if dataStr, ok := data.(string); ok {
				inputSize = len(dataStr)
			} else if dataSlice, ok := data.([]interface{}); ok {
				// If data is a slice, log its length
				inputSize = len(dataSlice)
			} else if dataMap, ok := data.(map[string]interface{}); ok {
				// If data is a map, log its length
				inputSize = len(dataMap)
			}
		} else {
			inputSize = len(input)
		}
	}
	
	logger.LogInfo("Executing connector", 
		"node_id", node.ID, 
		"connector_type", node.NodeType, 
		"input_size", inputSize)

	// Execute the connector
	result, err := connector.Execute(ctx, input)
	if err != nil {
		logger.LogError("Connector execution failed", 
			"node_id", node.ID, 
			"connector_type", node.NodeType, 
			"error", err.Error())
		return nil, fmt.Errorf("connector execution failed for node %s: %w", node.ID, err)
	}

	// Log result data size for debugging
	var resultSize int
	if result != nil {
		if data, ok := result["data"]; ok {
			// If data is a string, log its length
			if dataStr, ok := data.(string); ok {
				resultSize = len(dataStr)
			} else if dataSlice, ok := data.([]interface{}); ok {
				// If data is a slice, log its length
				resultSize = len(dataSlice)
			} else if dataMap, ok := data.(map[string]interface{}); ok {
				// If data is a map, log its length
				resultSize = len(dataMap)
			}
		} else {
			resultSize = len(result)
		}
	}
	
	logger.LogInfo("Connector execution completed", 
		"node_id", node.ID, 
		"connector_type", node.NodeType, 
		"result_size", resultSize)

	return result, nil
}

func (e *WorkflowEngine) updateExecutionStatus(
	executionID string,
	status string,
	errorMessage string,
	startTime time.Time,
	logs []map[string]interface{},
	results map[string]interface{},
) {
	execution, err := query.FindById[*wfModels.WorkflowExecution](executionID)
	if err != nil {
		logger.Error.Printf("Failed to find execution record %s: %v", executionID, err)
		return
	}

	endTime := time.Now()
	duration := int(endTime.Sub(startTime).Milliseconds())

	execution.Status = status
	execution.EndTime = pbTypes.NowDateTime()
	execution.Duration = duration

	if errorMessage != "" {
		execution.ErrorMessage = errorMessage
	}

	if logs != nil {
		logsJSON, err := json.Marshal(logs)
		if err == nil {
			execution.Logs = string(logsJSON)
		}
	}

	if results != nil {
		resultsJSON, err := json.Marshal(results)
		if err == nil {
			execution.Results = string(resultsJSON)
		}
	}

	// Save the record
	if err := query.SaveRecord(execution); err != nil {
		logger.Error.Printf("Failed to update execution record %s: %v", executionID, err)
	}
}

// updateExecutionLogs updates just the logs of a workflow execution
func (e *WorkflowEngine) updateExecutionLogs(executionID string, logs []map[string]interface{}) {
	execution, err := query.FindById[*wfModels.WorkflowExecution](executionID)
	if err != nil {
		logger.Error.Printf("Failed to find execution record %s: %v", executionID, err)
		return
	}

	logsJSON, err := json.Marshal(logs)
	if err != nil {
		logger.Error.Printf("Failed to marshal logs for execution %s: %v", executionID, err)
		return
	}

	execution.Logs = string(logsJSON)

	if err := query.SaveRecord(execution); err != nil {
		logger.Error.Printf("Failed to update logs for execution %s: %v", executionID, err)
	}
}

// RegisterConnectors registers all available connectors with the registry
// This is exported for use by the application to get available connectors
func RegisterConnectors(registry types.ConnectorRegistry) {
	// Register built-in connectors from the workflow package
	registry.Register("file_source", NewFileSourceConnector)
	registry.Register("transform_processor", NewTransformProcessor)
	registry.Register("log_destination", NewLogDestinationConnector)
	
	// Register connectors from the connectors package
	connectors.RegisterAllConnectors(registry)
	
	logger.Info.Println("All workflow connectors registered successfully")
}