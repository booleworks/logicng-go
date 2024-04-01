package normalform

import (
	f "booleworks.com/logicng/formula"
	"booleworks.com/logicng/handler"
)

// FactorizationHandler is a special handler able to abort CNF or DNF factorizations.
type FactorizationHandler interface {
	handler.Handler
	PerformedDistribution() bool
	CreatedClause(clause f.Formula) bool
}

// CNFHandler is a special handler for the advanced CNF algorithm.
type CNFHandler struct {
	aborted              bool
	distributionBoundary int
	clauseBoundary       int
	currentDistributions int
	currentClauses       int
}

// NewCNFHandler returns a new handler for the advanced CNF transformation with
// the given distribution and clause boundary.
func NewCNFHandler(distributionBoundary, clauseBoundary int) *CNFHandler {
	return &CNFHandler{
		false,
		distributionBoundary,
		clauseBoundary,
		0,
		0,
	}
}

// Started is called when the CNF factorization starts.
func (h *CNFHandler) Started() {
	h.aborted = false
	h.currentDistributions = 0
	h.currentClauses = 0
}

// Aborted reports whether the handler was aborted.
func (h *CNFHandler) Aborted() bool {
	return h.aborted
}

// PerformedDistribution is called each time a distribution during the
// factorization is performed and returns true if the computation should be
// continued.
func (h *CNFHandler) PerformedDistribution() bool {
	h.currentDistributions++
	h.aborted = h.distributionBoundary != -1 && h.currentDistributions > h.distributionBoundary
	return !h.aborted
}

// CreatedClause is called each time a clause is created during the
// factorization and returns true if the computation should be continued.
func (h *CNFHandler) CreatedClause(clause f.Formula) bool {
	h.currentClauses++
	h.aborted = h.clauseBoundary != -1 && h.currentClauses > h.clauseBoundary
	return !h.aborted
}
