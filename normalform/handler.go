package normalform

import "github.com/booleworks/logicng-go/event"

// FactorizationHandler is a handler for CNF and DNF factorization.
type FactorizationHandler struct {
	cancelled            bool
	distributionBoundary int
	clauseBoundary       int
	currentDistributions int
	currentClauses       int
}

// NewFactorizationHandler returns a new handler for the advanced CNF
// transformation with the given distribution and clause boundary.
func NewFactorizationHandler(distributionBoundary, clauseBoundary int) *FactorizationHandler {
	return &FactorizationHandler{false, distributionBoundary, clauseBoundary, 0, 0}
}

// ShouldResume processes the given event type and returns true if the
// computation should be resumed and false if it should be cancelled.
func (f *FactorizationHandler) ShouldResume(e event.Event) bool {
	if e == event.FactorizationStarted {
		f.currentDistributions = 0
		f.currentClauses = 0
	} else if e == event.DistributionPerformed {
		f.currentDistributions++
		f.cancelled = f.distributionBoundary != -1 && f.currentDistributions > f.distributionBoundary
	} else if e == event.FactorizationCreatedClause {
		f.currentClauses++
		f.cancelled = f.clauseBoundary != -1 && f.currentClauses > f.clauseBoundary
	}
	return !f.cancelled
}
