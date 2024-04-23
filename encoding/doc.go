// Package encoding provides data structures and algorithms for cardinality and
// pseudo-Boolean encodings in LogicNG.
//
// Encoders transform a constraint into a CNF representation.  For conventional
// (non-incremental) cardinality constraints we have the following encodings
// for AMO (at-most-one), EXO (exactly-one), AMK (at-most-k), ALK (at-least-k)
// constraints, and EXK (exactly-k):
//   - Pure (AMO, EXO): 'naive' encoding with no introduction of new variables
//     but quadratic size
//   - Ladder (AMO, EXO): Ladder/Regular Encoding due to Gent & Nightingale
//   - Product (AMO, EXO): the 2-Product Method due to Chen
//   - Nested (AMO, EXO): Nested pure encoding
//   - Commander (AMO, EXO): Commander Encoding due to Klieber & Kwon
//   - Binary (AMO, EXO): Binary Encoding due to Doggett, Frisch, Peugniez,
//     and Nightingale
//   - Bimander (AMO, EXO): Bimander Encoding due to Hölldobler and Nguyen
//   - Cardinality Network (AMK, ALK, EXK): Cardinality Network Encoding due to
//     Asín, Nieuwenhuis, Oliveras, and Rodríguez-Carbonell
//   - Totalizer (AMK, ALK, EXK): Totalizer Encoding due to Bailleux and
//     Boufkhad
//   - Modulo Totalizer (AMK, ALK): Modulo Totalizer due to Ogawa, Liu,
//     Hasegawa, Koshimura & Fujita
//
// Incremental cardinality constraints are a special variant of encodings,
// where the upper/lower bound of the constraint can be tightened by adding
// additional clauses to the resulting CNF.  All AMK and ALK encodings
// support incremental encodings.
//
// For pseudo-Boolean constraints there are three different encodings:
//   - PBAdderNetworks: Adder networks encoding
//   - PBSWC: A sequential weight counter for the encoding of pseudo-Boolean
//     constraints in CNF
//   - PBBinaryMerge: Binary merge due to Manthey, Philipp, and Steinke
//
// To encode a constraint explicitly (and not implicitly within e.g. the CNF
// methods) you can use the following code:
//
//	fac := formula.NewFactory()
//	cc := fac.AMO(fac.Vars("a", "b", "c", "d")...) // a + b + c + d <= 1
//	encoding, err := encoding.EncodeCC(fac, cc)
package encoding
