package workflow

import (
	"fmt"

	"github.com/shashank-sharma/backend/internal/services/workflow/types"
)

// ConnectorRegistryImpl implements the ConnectorRegistry interface
type ConnectorRegistryImpl struct {
	connectors map[string]types.ConnectorConstructor
}

// NewConnectorRegistryImpl creates a new connector registry
func NewConnectorRegistryImpl() *ConnectorRegistryImpl {
	registry := &ConnectorRegistryImpl{
		connectors: make(map[string]types.ConnectorConstructor),
	}
	
	return registry
}

// Register adds a connector to the registry
func (r *ConnectorRegistryImpl) Register(id string, constructor types.ConnectorConstructor) {
	r.connectors[id] = constructor
}

// Get returns a connector instance by ID, or nil if not found
func (r *ConnectorRegistryImpl) Get(id string) types.Connector {
	constructor, exists := r.connectors[id]
	if !exists {
		return nil
	}
	return constructor()
}

// Create instantiates a new connector of the specified type
func (r *ConnectorRegistryImpl) Create(id string) (types.Connector, error) {
	factory, exists := r.connectors[id]
	if !exists {
		return nil, fmt.Errorf("connector type not found: %s", id)
	}
	
	return factory(), nil
}

// GetAvailableConnectors returns a list of all registered connector types
func (r *ConnectorRegistryImpl) GetAvailableConnectors() []string {
	result := make([]string, 0, len(r.connectors))
	for connType := range r.connectors {
		result = append(result, connType)
	}
	return result
} 