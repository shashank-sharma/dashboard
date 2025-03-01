package workflow

import (
	"github.com/shashank-sharma/backend/internal/services/workflow/types"
)

// ConnectorType is an alias to types.ConnectorType for backward compatibility
type ConnectorType = types.ConnectorType

// Define connector type constants for backward compatibility
const (
	SourceConnector      = types.SourceConnector
	DestinationConnector = types.DestinationConnector
	ProcessorConnector   = types.ProcessorConnector
)

// Connector is an alias to types.Connector for backward compatibility
type Connector = types.Connector

// The deprecated BaseConnector has been removed
// Use types.BaseConnector directly instead 