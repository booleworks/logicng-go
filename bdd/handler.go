package bdd

import "booleworks.com/logicng/handler"

// A Handler for a BDD can abort the compilation of a BDD.  The method
// NewRefAdded is called by the BDD compiler everytime a newBdd BDD node
// reference is added.
type Handler interface {
	handler.Handler
	NewRefAdded() bool
}

// A TimeoutHandler aborts the BDD compilation dependent on the time it takes
// to generate the BDD.
type TimeoutHandler struct {
	handler.Timeout
}

// HandlerWithTimeout returns a newBdd BDD TimoutHandler for the given timout.
func HandlerWithTimeout(timeout handler.Timeout) *TimeoutHandler {
	return &TimeoutHandler{timeout}
}

// NewRefAdded is called by the BDD compiler everytime a newBdd BDD node reference
// is added.
func (t *TimeoutHandler) NewRefAdded() bool {
	return !t.TimeLimitExceeded()
}

// A NodesHandler aborts the BDD compilation dependent on the number of nodes
// which are generated during the compilation.
type NodesHandler struct {
	handler.Computation
	bound int
	count int
}

// HandlerWithNodes returns a newBdd BDD NodesHandler for the given bound of BDD
// nodes.
func HandlerWithNodes(bound int) *NodesHandler {
	return &NodesHandler{handler.Computation{}, bound, 0}
}

// NewRefAdded is called by the BDD compiler everytime a newBdd BDD node reference
// is added.
func (n *NodesHandler) NewRefAdded() bool {
	n.count++
	n.SetAborted(n.count >= n.bound)
	return !n.Aborted()
}
