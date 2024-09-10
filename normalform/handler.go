package normalform

import "github.com/booleworks/logicng-go/event"

// FactorizationHandler is a handler for CNF and DNF factorization.
type FactorizationHandler struct {
	canceled             bool
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
// computation should be resumed and false if it should be canceled.
func (f *FactorizationHandler) ShouldResume(e event.Event) bool {
	switch e {
	case event.FactorizationStarted:
		f.currentDistributions = 0
		f.currentClauses = 0
	case event.DistributionPerformed:
		f.currentDistributions++
		f.canceled = f.distributionBoundary != -1 && f.currentDistributions > f.distributionBoundary
	case event.FactorizationCreatedClause:
		f.currentClauses++
		f.canceled = f.clauseBoundary != -1 && f.currentClauses > f.clauseBoundary
	}
	return !f.canceled
}
