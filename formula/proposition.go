package formula

import "fmt"

// A Proposition is a formula with additional information like a textual
// description or a user-provided object.
type Proposition interface {
	Formula() Formula
	String() string
	Sprint(fac Factory) string
}

// A StandardProposition is a simple proposition implementation with a formula
// and a textual description.
type StandardProposition struct {
	form        Formula
	Description string
}

// An ExtendedProposition is a proposition implementation with a formula and a
// user provided generic backpack.
type ExtendedProposition[T fmt.Stringer] struct {
	form     Formula
	Backpack T
}

// Formula returns the formula of the proposition.
func (p *StandardProposition) Formula() Formula { return p.form }

func (p *StandardProposition) String() string {
	return fmt.Sprintf("%s: %s", p.Description, p.form)
}

// Sprint takes a formula factory and pretty-prints the proposition.
func (p *StandardProposition) Sprint(fac Factory) string {
	return fmt.Sprintf("%s: %s", p.Description, p.form.Sprint(fac))
}

// Formula returns the formula of the proposition.
func (p *ExtendedProposition[T]) Formula() Formula { return p.form }

func (p *ExtendedProposition[T]) String() string {
	return fmt.Sprintf("%s: %s", p.Backpack, p.form)
}

// Sprint takes a formula factory and pretty-prints the proposition.
func (p *ExtendedProposition[T]) Sprint(fac Factory) string {
	return fmt.Sprintf("%s: %s", p.Backpack, p.form.Sprint(fac))
}

// NewStandardProposition returns a new standard proposition with the formula
// and an optional textual description.
func NewStandardProposition(formula Formula, description ...string) *StandardProposition {
	var desc string
	if len(description) > 0 {
		desc = description[0]
	}
	return &StandardProposition{formula, desc}
}

// NewExtendedProposition returns a new extended proposition with the formula
// and the given backpack.
func NewExtendedProposition[T fmt.Stringer](formula Formula, backpack T) *ExtendedProposition[T] {
	return &ExtendedProposition[T]{formula, backpack}
}
