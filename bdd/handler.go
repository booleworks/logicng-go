package bdd

import (
	"github.com/booleworks/logicng-go/event"
)

// A NodesHandler cancels the BDD compilation dependent on the number of nodes
// which are generated during the compilation.
type NodesHandler struct {
	canceled bool
	bound    int
	count    int
}

// HandlerWithNodes returns a new BDD NodesHandler for the given bound of BDD
// nodes.
func HandlerWithNodes(bound int) *NodesHandler {
	return &NodesHandler{false, bound, 0}
}

// ShouldResume returns false if the bound of generated BDD nodes is reached.
func (n *NodesHandler) ShouldResume(e event.Event) bool {
	switch e {
	case event.BddComputationStarted:
		n.count = 0
	case event.BddNewRefAdded:
		n.count++
		n.canceled = n.count >= n.bound
	}
	return !n.canceled
}
