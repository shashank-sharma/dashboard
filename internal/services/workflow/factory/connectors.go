// Package factory provides factory functions for creating and registering workflow connectors
package factory

import (
	"github.com/shashank-sharma/backend/internal/services/workflow/types"
)

// RegisterAllConnectors registers all available connectors with the registry
// This is where we register connectors from other packages to avoid import cycles
func RegisterAllConnectors(registry types.ConnectorRegistry) {
	// Register connectors from the connectors package
	registerConnectorPackage(registry)
}

// Creates adapter functions to convert between connector types
// These help us avoid cyclic imports by using adapters
func wrapConstructor(constructor interface{}) types.ConnectorConstructor {
	// This is a simplification for illustration purposes
	// In reality, you would need to use reflection or implement type-specific adapters
	return func() types.Connector {
		// In a real implementation, call the original constructor and 
		// adapt the specific connector type to the types.Connector interface
		return nil
	}
}

// registerConnectorPackage registers connectors from the connectors package
func registerConnectorPackage(registry types.ConnectorRegistry) {
	// Since we cannot directly register connectors from different packages
	// due to type compatibility issues with interface definitions,
	// we would need to manually register each connector with proper adapter pattern
	
	// To avoid extensive code changes, consider:
	// 1. Define core connector interfaces in a common package
	// 2. Have both workflow and connectors packages import those common interfaces
	// 3. Update the registration logic to use those common interfaces
	
	// For now, this file serves as a placeholder for future connector registration
	// that avoids import cycles. The implementation requires a refactoring of the
	// connector interfaces across packages.
} 