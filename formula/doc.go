// Package formula provides all the main data-structures for working with
// propositional formulas in LogicNG and generate them via a formula factory.
//
// Usually the formula factory is the starting point of working with LogicNG.
//
//	fac := formula.NewFactory()
//
// From there on you can generate formulas on the factory or use a parser to
// parse them to the factory.  All major algorithms which generate or
// manipulate formulas will need a formula factory as parameter.  Usually
// within an application you only have one formula factory and keep working
// with it. Very seldom there should be the need to create more than one
// factory.
//
// Formula factories are not thread-safe and you can not share formulas between
// formula factories since a formula in fact is only a unique ID (unit32) on
// the factory.
package formula
