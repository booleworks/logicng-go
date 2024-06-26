// Package parser provides a parser for propositional and pseudo-Boolean
// formulas to LogicNG formula structures.  Variable names in LogicNG can begin
// with digits, therefore variable names like '300' are valid.  Variable names
// must begin with [A-Za-z0-9_@#] and can contain letters, digits, or [_#.].
// The LogicNG grammar for propositional formulas is:
//
//	constant true: `$true`
//	constant false: $false
//	negation: ~
//	implication: =>
//	equivalence: <=>
//	conjunction: &
//	disjunction: |
//	left parentheses: (
//	right parentheses: )
//
// For pseudo-Boolean formulas the following symbols can be used:
//
//	addition: +
//	minus: -
//	multiplication: *
//	equal: =
//	less-than: <
//	less-than or equal: <=
//	greater-than: >
//	greater-than or equal: >=
//
// There are two methods on a parser: Parse and ParseSafe.  Since often in use
// cases, you know, that the formula you are going to parse is syntactically
// correct (e.g. since it was generated by LogicNG in the first place) you can
// use the ParseUnsafe method which just panics when encountering a syntax
// error.  Usually however you should use the Parse method which provides a
// proper error value in case of syntax errors.
//
// Usage example of the parser:
//
//	fac := f.NewFactory()
//	parser := NewPropositionalParser(fac)
//	formula, err := parser.Parse("A & (B|C) => ~X") // proper error handling
//	formula = parser.ParseUnsafe("A & (B|C) => ~X") // panics on errors
package parser
