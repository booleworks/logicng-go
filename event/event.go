package event

type Event interface {
	EventType() string
}

type event struct {
	eventType string
}

func (e event) EventType() string {
	return e.eventType
}

var (
	FactorizationStarted          = event{"Factorization Started"}
	BddComputationStarted         = event{"BDD Computation Started"}
	DnnfComputationStarted        = event{"DNNF Computation Started"}
	SatCallStarted                = event{"SAT Call Started"}
	MaxSATCallStarted             = event{"Max-SAT Call Started"}
	BackboneComputationStarted    = event{"Backbone Computation Started"}
	AdvancedSimplificationStarted = event{"Advanced Simplification Started"}
	PrimeComputationStarted       = event{"Prime Computation Started"}
	ImplicantReductionStarted     = event{"Implicant Reduction Started"}
	ImplicateReductionStarted     = event{"Implicate Reduction Started"}
	MusComputationStarted         = event{"MUS Computation Started"}
	SmusComputationStarted        = event{"SMUS Computation Started"}
	OptimizationFunctionStarted   = event{"Optimization Function Started"}
	ModelEnumerationStarted       = event{"Model Enumeration Started"}

	SatCallFinished    = event{"SAT Call Finished"}
	MaxSatCallFinished = event{"Max-SAT Call Finished"}

	ModelEnumerationCommit              = event{"Model Enumeration Commit"}
	ModelEnumerationRollback            = event{"Model Enumeration Rollback"}
	FactorizationCreatedClause          = event{"Factorization Created Clause"}
	DistributionPerformed               = event{"Distribution Performed"}
	BddNewRefAdded                      = event{"BDD New Ref Added"}
	DnnfShannonExpansion                = event{"DNNF Shannon Expansion"}
	DnnfDtreeMinFillGraphInitialized    = event{"DNNF DTree MinFill Graph initialized"}
	DnnfDtreeMinFillNewIteration        = event{"DNNF DTree MinFill new iteration"}
	DnnfDtreeProcessingNextOrderVar     = event{"DNNF DTree processing next order variable"}
	SatConflictDetected                 = event{"SAT Conflict Detected"}
	SubsumptionStartingUbTreeGeneration = event{"Subsumption Starting UB Tree Generation"}
	SubsumptionAddedNewSet              = event{"Subsumption Added New Set"}

	Nothing = event{"Nothing"}
)
