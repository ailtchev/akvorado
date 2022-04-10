// Package broker synchronizes the different internal services.
package broker

import (
	"sync"

	"akvorado/common/http"
	"akvorado/common/reporter"
)

// Component represents the broker.
type Component struct {
	r      *reporter.Reporter
	d      *Dependencies
	config Configuration

	serviceLock           sync.Mutex
	serviceConfigurations map[ServiceType]interface{}
	registeredServices    map[ServiceType]map[string]bool
}

// Dependencies define the dependencies of the broker.
type Dependencies struct {
	HTTP *http.Component
}

// ServiceType describes the different internal services
type ServiceType string

var (
	// InletService represents the inlet service type
	InletService ServiceType = "inlet"
	// OrchestratorService represents the orchestrator service type
	OrchestratorService ServiceType = "orchestrator"
	// ConsoleService represents the console service type
	ConsoleService ServiceType = "console"
)

// New creates a new broker component.
func New(r *reporter.Reporter, configuration Configuration, dependencies Dependencies) (*Component, error) {
	c := Component{
		r:      r,
		d:      &dependencies,
		config: configuration,

		serviceConfigurations: map[ServiceType]interface{}{},
		registeredServices:    map[ServiceType]map[string]bool{},
	}

	c.d.HTTP.GinRouter.GET("/api/v0/orchestrator/broker/configuration/:service", c.configurationHandlerFunc)

	return &c, nil
}

// RegisterConfiguration registers the configuration for a service.
func (c *Component) RegisterConfiguration(service ServiceType, configuration interface{}) {
	c.serviceLock.Lock()
	c.serviceConfigurations[service] = configuration
	c.serviceLock.Unlock()
}
