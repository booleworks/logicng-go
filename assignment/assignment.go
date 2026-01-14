package assignment

import (
	"slices"
	"strings"

	"github.com/booleworks/logicng-go/errorx"
	f "github.com/booleworks/logicng-go/formula"
)

// An Assignment represents a mapping from Boolean variables to truth values
type Assignment struct {
	pos map[uint32]present
	neg map[uint32]present
}

// Empty creates an empty Assignment
func Empty() *Assignment {
	return &Assignment{make(map[uint32]present), make(map[uint32]present)}
}

// New creates a new Assignment with the given literals.  Literals with a
// positive phase are added as a mapping to true, literals with a negative
// phase as a mapping to false.  Returns an error, if the literals list
// contains complementary literals.
func New(fac f.Factory, literals ...f.Literal) (*Assignment, error) {
	ass := &Assignment{make(map[uint32]present), make(map[uint32]present)}
	for _, lit := range literals {
		err := ass.AddLit(fac, lit)
		if err != nil {
			return nil, err
		}
	}
	return ass, nil
}

// AddLit add a single literal to the assignment.  A literal with a
// positive phase is added as a mapping to true, a literal with a negative
// phase as a mapping to false.
func (a *Assignment) AddLit(fac f.Factory, literal f.Literal) error {
	if literal.IsPos() {
		if _, ok := a.neg[literal.ID()^1]; ok {
			return errorx.BadInput("%s (opposite phase present)", literal.Sprint(fac))
		}
		a.pos[literal.ID()] = present{}
	} else {
		if _, ok := a.pos[literal.ID()^1]; ok {
			return errorx.BadInput("%s (opposite phase present)", literal.Sprint(fac))
		}
		a.neg[literal.ID()] = present{}
	}
	return nil
}

// PosVars returns all variables of the assignment mapped to true.
func (a *Assignment) PosVars() []f.Variable {
	slice := make([]f.Variable, 0, len(a.pos))
	for l := range a.pos {
		slice = append(slice, f.EncodeVariable(l))
	}
	return slice
}

// NegVars returns all variables of the assignment mapped to false.
// Note, this function return the variables, not the negative literals.
func (a *Assignment) NegVars() []f.Variable {
	slice := make([]f.Variable, 0, len(a.neg))
	for l := range a.neg {
		slice = append(slice, f.EncodeVariable(l))
	}
	return slice
}

// Size returns the number of variables in the Assignment.
func (a *Assignment) Size() int {
	return len(a.pos) + len(a.neg)
}

// Sprint prints the assignment in human-readable form with deterministic
// ordering (sorted by variable name).
func (a *Assignment) Sprint(fac f.Factory) string {
	// Collect all literals
	literals := make([]f.Literal, 0, len(a.pos)+len(a.neg))
	for l := range a.pos {
		literals = append(literals, f.Literal(f.EncodeFormula(f.SortLiteral, l)))
	}
	for l := range a.neg {
		literals = append(literals, f.Literal(f.EncodeFormula(f.SortLiteral, l)))
	}

	// Sort by variable name for deterministic output
	slices.SortFunc(literals, func(a, b f.Literal) int {
		nameA, _ := fac.VarName(a.Variable())
		nameB, _ := fac.VarName(b.Variable())
		return strings.Compare(nameA, nameB)
	})

	// Build the output string
	var sb strings.Builder
	sb.WriteString("[")
	for i, lit := range literals {
		sb.WriteString(lit.Sprint(fac))
		if i < len(literals)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("]")
	return sb.String()
}

func (a *Assignment) evaluateVariable(variable f.Variable) bool {
	_, ok := a.pos[variable.ID()]
	return ok
}

func (a *Assignment) evaluateNegativeLiteral(literal f.Literal) bool {
	_, ok := a.neg[literal.ID()]
	if ok {
		return true
	}
	_, ok = a.pos[literal.ID()^1]
	return !ok
}

func (a *Assignment) restrictVariable(fac f.Factory, variable f.Variable) f.Formula {
	_, ok := a.pos[variable.ID()]
	if ok {
		return fac.Verum()
	}
	_, ok = a.neg[variable.ID()^1]
	if ok {
		return fac.Falsum()
	}
	return variable.AsFormula()
}

func (a *Assignment) restrictNegativeLiteral(fac f.Factory, literal f.Literal) f.Formula {
	_, ok := a.pos[literal.ID()^1]
	if ok {
		return fac.Falsum()
	}
	_, ok = a.neg[literal.ID()]
	if ok {
		return fac.Verum()
	}
	return literal.AsFormula()
}

type present struct{}
