package iter

import (
	"booleworks.com/logicng/configuration"
)

// Config represents a model iteration configuration including a handler
// and an iteration strategy.
type Config struct {
	Handler  Handler
	Strategy Strategy
}

// Sort returns the configuration sort (ModelIteration).
func (Config) Sort() configuration.Sort {
	return configuration.ModelIteration
}

// DefaultConfig returns the default configuration for a model iteration
// configuration.
func (Config) DefaultConfig() configuration.Config {
	return DefaultConfig()
}

// DefaultConfig returns the default configuration for a
// model iteration configuration.
func DefaultConfig() *Config {
	return &Config{nil, DefaultStrategy()}
}
