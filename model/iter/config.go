package iter

import (
	"github.com/booleworks/logicng-go/configuration"
	"github.com/booleworks/logicng-go/handler"
)

// Config represents a model iteration configuration including a handler
// and an iteration strategy.
type Config struct {
	Handler  handler.Handler
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
	return &Config{handler.NopHandler, DefaultStrategy()}
}
