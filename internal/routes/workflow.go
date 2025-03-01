package routes

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/pocketbase/pocketbase/core"
	"github.com/shashank-sharma/backend/internal/logger"
	"github.com/shashank-sharma/backend/internal/models"
	"github.com/shashank-sharma/backend/internal/query"
	"github.com/shashank-sharma/backend/internal/services/workflow"
	"github.com/shashank-sharma/backend/internal/services/workflow/connectors"
	"github.com/shashank-sharma/backend/internal/util"
)

// RegisterWorkflowRoutes registers the workflow-related API endpoints
func RegisterWorkflowRoutes(e *core.ServeEvent, engine *workflow.WorkflowEngine) {
	logger.LogInfo("Registering workflow routes")
	// Create a new registry
	connRegistry := workflow.NewConnectorRegistryImpl()
	
	// Register connectors directly
	connectors.RegisterAllConnectors(connRegistry)
	workflow.RegisterLocalConnectors(connRegistry)
	
	// Get available connectors
	e.Router.GET("/api/workflows/connectors", func(e *core.RequestEvent) error {
		// Get connector IDs
		connectorIDs := connRegistry.GetAvailableConnectors()
		
		// Convert to richer metadata format
		results := make([]map[string]interface{}, 0, len(connectorIDs))
		
		for _, id := range connectorIDs {
			connector := connRegistry.Get(id)
			if connector != nil {
				results = append(results, map[string]interface{}{
					"id":           connector.ID(),
					"name":         connector.Name(),
					"type":         connector.Type(),
					"configSchema": connector.GetConfigSchema(),
				})
			}
		}
		
		return e.JSON(http.StatusOK, map[string]interface{}{
			"connectors": results,
		})
	})

	e.Router.POST("/api/workflows/{id}/execute", func(e *core.RequestEvent) error {
		workflowId := e.Request.PathValue("id")
		if workflowId == "" {
			return e.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "Missing workflow ID",
			})
		}

		token := e.Request.Header.Get("Authorization")
		userId, err := util.GetUserId(token)
		if err != nil {
			return e.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Unauthorized"})
		}

		workflow, err := query.FindByFilter[*models.Workflow](map[string]interface{}{
			"id": workflowId,
			"user": userId,
		})
		if err != nil {
			return e.JSON(http.StatusNotFound, map[string]interface{}{
				"error": "Workflow not found: " + err.Error(),
			})
		}

		ctx := context.WithValue(context.Background(), "user", userId)

		// Execute the workflow with user context
		execution, err := engine.ExecuteWorkflow(ctx, workflow.Id)
		if err != nil {
			return e.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "Failed to execute workflow: " + err.Error(),
			})
		}

		return e.JSON(http.StatusAccepted, execution)
	})

	// Get execution status
	e.Router.GET("/api/workflows/executions/{id}", func(e *core.RequestEvent) error {
		executionId := e.Request.PathValue("id")
		if executionId == "" {
			return e.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "Missing execution ID",
			})
		}

		// Get user ID from token
		token := e.Request.Header.Get("Authorization")
		userId, err := util.GetUserId(token)
		if err != nil {
			return e.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Unauthorized"})
		}

		// Get the execution record to check access
		execution, err := query.FindByFilter[*models.WorkflowExecution](map[string]interface{}{
			"id": executionId,
		})
		if err != nil {
			return e.JSON(http.StatusNotFound, map[string]interface{}{
				"error": "Execution not found",
			})
		}

		// Get the associated workflow to check access
		_, err = query.FindByFilter[*models.Workflow](map[string]interface{}{
			"id": execution.WorkflowID,
			"user": userId,
		})
		if err != nil {
			return e.JSON(http.StatusNotFound, map[string]interface{}{
				"error": "Associated workflow not found or you don't have access",
			})
		}

		// Get the execution status
		status, err := engine.GetExecutionStatus(executionId)
		if err != nil {
			return e.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "Failed to get execution status: " + err.Error(),
			})
		}

		return e.JSON(http.StatusOK, status)
	})

	// Register webhook triggers for workflows
	e.Router.POST("/api/workflows/{id}/webhook", func(e *core.RequestEvent) error {
		workflowId := e.Request.PathValue("id")
		if workflowId == "" {
			return e.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "Missing workflow ID",
			})
		}

		// Get user ID from token
		token := e.Request.Header.Get("Authorization")
		userId, err := util.GetUserId(token)
		if err != nil {
			return e.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Unauthorized"})
		}

		// Load the workflow to check if it exists and user has access
		workflow, err := query.FindByFilter[*models.Workflow](map[string]interface{}{
			"id": workflowId,
			"user": userId,
		})
		if err != nil {
			return e.JSON(http.StatusNotFound, map[string]interface{}{
				"error": "Workflow not found or unauthorized",
			})
		}

		// Parse the webhook payload
		var payload map[string]interface{}
		decoder := json.NewDecoder(e.Request.Body)
		if err := decoder.Decode(&payload); err != nil {
			return e.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "Invalid payload format",
			})
		}

		// Create context with user ID
		ctx := context.WithValue(context.Background(), "user", userId)

		// Execute the workflow with the webhook payload and user context
		execution, err := engine.ExecuteWorkflow(ctx, workflow.Id)
		if err != nil {
			return e.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "Failed to execute workflow: " + err.Error(),
			})
		}

		return e.JSON(http.StatusAccepted, execution)
	})
}

// The helper function is no longer needed as it's been moved to the workflow package 