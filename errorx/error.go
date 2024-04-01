package errorx

import (
	"errors"
	"fmt"
)

var (
	ErrBadFormulaSort = errors.New("bad formula type")   // formula type is not allowed here
	ErrBadInput       = errors.New("bad input")          // input is not allowed here
	ErrUnknownFormula = errors.New("unknown formula")    // formula is not known with this sort on the factory
	ErrUnknownEnumVal = errors.New("unknown enum value") // enum value is unknown
	ErrIllegalState   = errors.New("illegal state")      // an illegal internal state has been reached
)

// BadFormulaSort returns an error for an unsupported formula sort.
func BadFormulaSort(sort fmt.Stringer) error {
	return fmt.Errorf("%w: %s", ErrBadFormulaSort, sort)
}

// BadInput returns an error for an unsupported input.
func BadInput(text string, params ...any) error {
	return fmt.Errorf("%w: %s", ErrBadInput, fmt.Sprintf(text, params...))
}

// UnknownFormula returns an error for an unknown formula.
func UnknownFormula(formula fmt.Stringer) error {
	return fmt.Errorf("%w: %s", ErrUnknownFormula, formula)
}

// UnknownEnumValue returns an error for an unknown enum value.
func UnknownEnumValue(value fmt.Stringer) error {
	return fmt.Errorf("%w: %s", ErrUnknownEnumVal, value)
}

// IllegalState returns an error for an illegal internal state.
func IllegalState(text string, params ...any) error {
	return fmt.Errorf("%w: %s", ErrIllegalState, fmt.Sprintf(text, params...))
}
