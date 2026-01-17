package model

import (
	"strings"

	"github.com/booleworks/logicng-go/assignment"
	f "github.com/booleworks/logicng-go/formula"
)

// A Model represents a list of literals.
type Model struct {
	Literals []f.Literal
}

// New returns a new model with the given literals
func New(literals ...f.Literal) *Model {
	if literals == nil {
		return &Model{[]f.Literal{}}
	}
	return &Model{literals}
}

// AddLiteral adds the given literals to the model.
func (m *Model) AddLiteral(literals ...f.Literal) {
	m.Literals = append(m.Literals, literals...)
}

// FromAssignment generates a new model from a given assignment.
func FromAssignment(fac f.Factory, ass assignment.Assignment) *Model {
	literals := make([]f.Literal, 0, ass.Size())
	literals = append(literals, f.VariablesAsLiterals(ass.PosVars())...)
	for _, lit := range ass.NegVars() {
		literals = append(literals, lit.Negate(fac))
	}
	return &Model{literals}
}

// Size returns the size of the model.
func (m *Model) Size() int {
	return len(m.Literals)
}

// Assignment returns an assignment from the model.  Returns an error if this
// model contains complementary literals (which should not happen if used
// right).
func (m *Model) Assignment(fac f.Factory) (*assignment.Assignment, error) {
	return assignment.New(fac, m.Literals...)
}

// Formula returns a formula for the model which is the conjunction of its literals.
func (m *Model) Formula(fac f.Factory) f.Formula {
	return fac.Minterm(m.Literals...)
}

// PosVars returns the positive variables of the model.
func (m *Model) PosVars() []f.Variable {
	vars := make([]f.Variable, 0, len(m.Literals)/2)
	for _, lit := range m.Literals {
		if lit.IsPos() {
			vars = append(vars, f.Variable(lit))
		}
	}
	return vars
}

// NegLits returns the negative literals of the model.
func (m *Model) NegLits() []f.Literal {
	vars := make([]f.Literal, 0, len(m.Literals)/2)
	for _, lit := range m.Literals {
		if lit.IsNeg() {
			vars = append(vars, lit)
		}
	}
	return vars
}

// NegVars returns the negative variables of the model.
func (m *Model) NegVars() []f.Variable {
	vars := make([]f.Variable, 0, len(m.Literals)/2)
	for _, lit := range m.Literals {
		if lit.IsNeg() {
			variable := lit.Variable()
			vars = append(vars, variable)
		}
	}
	return vars
}

// Sprint prints the model in human-readable form.
func (m *Model) Sprint(fac f.Factory) string {
	var sb strings.Builder
	sb.WriteString("[")
	length := len(m.Literals)
	for i, lit := range m.Literals {
		sb.WriteString(lit.Sprint(fac))
		if i+1 < length {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("]")
	return sb.String()
}
