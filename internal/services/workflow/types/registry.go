package types

// ConnectorConstructor is a function that creates a new connector instance
type ConnectorConstructor func() Connector

// ConnectorRegistry is an interface for managing connector registrations
type ConnectorRegistry interface {
	// Register registers a connector constructor with the registry
	Register(id string, constructor ConnectorConstructor)
	
	// Get returns a connector instance by ID, or nil if not found
	Get(id string) Connector
	
	// Create instantiates a new connector of the specified type
	Create(id string) (Connector, error)
	
	// GetAvailableConnectors returns a list of all registered connector IDs
	GetAvailableConnectors() []string
} 