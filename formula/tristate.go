package formula

// Tristate represents a Boolean value with the possibility of a third UNDEF
// value.
type Tristate byte

const (
	TristateTrue Tristate = iota
	TristateFalse
	TristateUndef
)

//go:generate stringer -type=Tristate

// Negate returns the negation of the tristate TRUE turns FALSE, FALSE turns
// TRUE, and UNDEF stays UNDEF.
func (t Tristate) Negate() Tristate {
	switch t {
	case TristateTrue:
		return TristateFalse
	case TristateFalse:
		return TristateTrue
	default:
		return TristateUndef
	}
}

// TristateFromBool returns the tristate for the given Boolean value.
func TristateFromBool(b bool) Tristate {
	if b {
		return TristateTrue
	}
	return TristateFalse
}
