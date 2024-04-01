package formula

type testdata struct {
	F     Factory
	True  Formula
	False Formula
	A     Formula
	B     Formula
	C     Formula
	D     Formula
	X     Formula
	Y     Formula
	NA    Formula
	NB    Formula
	NX    Formula
	NY    Formula
	VA    Variable
	VB    Variable
	VC    Variable
	VD    Variable
	VX    Variable
	VY    Variable
	LA    Literal
	LB    Literal
	LC    Literal
	LD    Literal
	LX    Literal
	LY    Literal
	LNA   Literal
	LNB   Literal
	LNX   Literal
	LNY   Literal
	OR1   Formula
	OR2   Formula
	OR3   Formula
	AND1  Formula
	AND2  Formula
	AND3  Formula
	NOT1  Formula
	NOT2  Formula
	IMP1  Formula
	IMP2  Formula
	IMP3  Formula
	IMP4  Formula
	EQ1   Formula
	EQ2   Formula
	EQ3   Formula
	EQ4   Formula
	PBC1  Formula
	PBC2  Formula
	PBC3  Formula
	PBC4  Formula
	PBC5  Formula
	PB1   Formula
	PB2   Formula
	CC1   Formula
	CC2   Formula
	AMO1  Formula
	AMO2  Formula
	EXO1  Formula
	EXO2  Formula
}

func NewTestData(fac Factory) testdata {
	data := testdata{
		F:     fac,
		True:  fac.Verum(),
		False: fac.Falsum(),
		A:     fac.Variable("a"),
		B:     fac.Variable("b"),
		C:     fac.Variable("c"),
		D:     fac.Variable("d"),
		X:     fac.Variable("x"),
		Y:     fac.Variable("y"),
		NA:    fac.Literal("a", false),
		NB:    fac.Literal("b", false),
		NX:    fac.Literal("x", false),
		NY:    fac.Literal("y", false),
		VA:    fac.Var("a"),
		VB:    fac.Var("b"),
		VC:    fac.Var("c"),
		VD:    fac.Var("d"),
		VX:    fac.Var("x"),
		VY:    fac.Var("y"),
		LA:    fac.Lit("a", true),
		LB:    fac.Lit("b", true),
		LC:    fac.Lit("c", true),
		LD:    fac.Lit("d", true),
		LX:    fac.Lit("x", true),
		LY:    fac.Lit("y", true),
		LNA:   fac.Lit("a", false),
		LNB:   fac.Lit("b", false),
		LNX:   fac.Lit("x", false),
		LNY:   fac.Lit("y", false),
	}

	data.OR1 = fac.Or(data.X, data.Y)
	data.OR2 = fac.Or(data.NX, data.NY)
	data.OR3 = fac.Or(fac.And(data.A, data.B), fac.And(data.NA, data.NB))
	data.AND1 = fac.And(data.A, data.B)
	data.AND2 = fac.And(data.NA, data.NB)
	data.AND3 = fac.And(data.OR1, data.OR2)
	data.NOT1 = fac.Not(data.AND1)
	data.NOT2 = fac.Not(data.OR1)
	data.IMP1 = fac.Implication(data.A, data.B)
	data.IMP2 = fac.Implication(data.NA, data.NB)
	data.IMP3 = fac.Implication(data.AND1, data.OR1)
	data.IMP4 = fac.Implication(fac.Equivalence(data.A, data.B), fac.Equivalence(data.NX, data.NY))
	data.EQ1 = fac.Equivalence(data.A, data.B)
	data.EQ2 = fac.Equivalence(data.NA, data.NB)
	data.EQ3 = fac.Equivalence(data.AND1, data.OR1)
	data.EQ4 = fac.Equivalence(data.IMP1, data.IMP2)

	literals := []Literal{Literal(data.A), Literal(data.B), Literal(data.X)}
	coefficients := []int{2, -4, 3}
	lits1 := []Variable{Variable(data.A)}
	lits2 := []Literal{Literal(data.A), Literal(data.NB), Literal(data.C)}
	litsCc2 := []Variable{Variable(data.A), Variable(data.B), Variable(data.C)}
	coeffs1 := []int{3}
	coeffs2 := []int{3, -2, 7}

	data.PBC5 = fac.PBC(LE, 2, literals, coefficients)
	data.PBC1 = fac.PBC(EQ, 2, literals, coefficients)
	data.PBC2 = fac.PBC(GT, 2, literals, coefficients)
	data.PBC3 = fac.PBC(GE, 2, literals, coefficients)
	data.PBC4 = fac.PBC(LT, 2, literals, coefficients)
	data.PB1 = fac.PBC(LE, 2, VariablesAsLiterals(lits1), coeffs1)
	data.PB2 = fac.PBC(LE, 8, lits2, coeffs2)
	data.CC1 = fac.CC(LT, 1, lits1...)
	data.CC2 = fac.CC(GE, 2, litsCc2...)
	data.AMO1 = fac.AMO(lits1...)
	data.AMO2 = fac.AMO(litsCc2...)
	data.EXO1 = fac.EXO(lits1...)
	data.EXO2 = fac.EXO(litsCc2...)
	return data
}

func NewCornerCases(fac Factory) []Formula {
	formulas := make([]Formula, 8)
	formulas[0] = fac.Falsum()
	formulas[1] = fac.Not(fac.Falsum())
	formulas[2] = fac.Verum()
	formulas[3] = fac.Not(fac.Verum())
	formulas[4] = fac.Variable("a")
	formulas[5] = fac.Literal("a", false)
	formulas[6] = fac.Not(fac.Variable("a"))
	formulas[7] = fac.Not(fac.Not(fac.Not(fac.Variable("a"))))
	formulas = append(formulas, binaryCornerCases(SortImpl, fac)...)
	formulas = append(formulas, binaryCornerCases(SortEquiv, fac)...)
	formulas = append(formulas, naryCornerCases(SortOr, fac)...)
	formulas = append(formulas, naryCornerCases(SortAnd, fac)...)
	formulas = append(formulas, cornerCasesPBC(fac)...)
	return formulas
}

func binaryCornerCases(sort FSort, fac Factory) []Formula {
	formulas := make([]Formula, 20)
	a := fac.Variable("a")
	na := fac.Literal("a", false)
	b := fac.Variable("b")
	nb := fac.Literal("b", false)

	formulas[0], _ = fac.BinaryOperator(sort, fac.Verum(), fac.Verum())
	formulas[1], _ = fac.BinaryOperator(sort, fac.Falsum(), fac.Verum())
	formulas[2], _ = fac.BinaryOperator(sort, fac.Verum(), fac.Falsum())
	formulas[3], _ = fac.BinaryOperator(sort, fac.Falsum(), fac.Falsum())

	formulas[4], _ = fac.BinaryOperator(sort, fac.Verum(), a)
	formulas[5], _ = fac.BinaryOperator(sort, a, fac.Verum())
	formulas[6], _ = fac.BinaryOperator(sort, fac.Verum(), na)
	formulas[7], _ = fac.BinaryOperator(sort, na, fac.Verum())

	formulas[8], _ = fac.BinaryOperator(sort, fac.Falsum(), a)
	formulas[9], _ = fac.BinaryOperator(sort, a, fac.Falsum())
	formulas[10], _ = fac.BinaryOperator(sort, fac.Falsum(), na)
	formulas[11], _ = fac.BinaryOperator(sort, na, fac.Falsum())

	formulas[12], _ = fac.BinaryOperator(sort, a, a)
	formulas[13], _ = fac.BinaryOperator(sort, a, na)
	formulas[14], _ = fac.BinaryOperator(sort, na, a)
	formulas[15], _ = fac.BinaryOperator(sort, na, na)

	formulas[16], _ = fac.BinaryOperator(sort, a, b)
	formulas[17], _ = fac.BinaryOperator(sort, a, nb)
	formulas[18], _ = fac.BinaryOperator(sort, na, b)
	formulas[19], _ = fac.BinaryOperator(sort, na, nb)

	return formulas
}

func naryCornerCases(sort FSort, fac Factory) []Formula {
	formulas := make([]Formula, 15)
	a := fac.Variable("a")
	na := fac.Literal("a", false)
	b := fac.Variable("b")
	nb := fac.Literal("b", false)
	c := fac.Variable("c")
	nc := fac.Literal("c", false)

	formulas[0], _ = fac.NaryOperator(sort)

	formulas[1], _ = fac.NaryOperator(sort, fac.Falsum())
	formulas[2], _ = fac.NaryOperator(sort, fac.Verum())
	formulas[3], _ = fac.NaryOperator(sort, fac.Falsum(), fac.Verum())

	formulas[4], _ = fac.NaryOperator(sort, a)
	formulas[5], _ = fac.NaryOperator(sort, na)

	formulas[6], _ = fac.NaryOperator(sort, fac.Verum(), a)
	formulas[7], _ = fac.NaryOperator(sort, fac.Verum(), na)
	formulas[8], _ = fac.NaryOperator(sort, fac.Falsum(), a)
	formulas[9], _ = fac.NaryOperator(sort, fac.Falsum(), na)

	formulas[10], _ = fac.NaryOperator(sort, a, na)
	formulas[11], _ = fac.NaryOperator(sort, a, b)
	formulas[12], _ = fac.NaryOperator(sort, a, b, c)
	formulas[13], _ = fac.NaryOperator(sort, na, nb, nc)
	formulas[14], _ = fac.NaryOperator(sort, a, b, c, na)
	return formulas
}

func cornerCasesPBC(fac Factory) []Formula {
	formulas := make([]Formula, 0)
	for _, sort := range []CSort{LE, LT, GT, GE, EQ} {
		formulas = append(formulas, cornerCasePBC2(fac, sort)...)
	}
	return formulas
}

func cornerCasePBC2(fac Factory, comp CSort) []Formula {
	formulas := make([]Formula, 0)
	a := fac.Lit("a", true)
	na := fac.Lit("a", false)
	b := fac.Lit("b", true)
	nb := fac.Lit("b", false)
	c := fac.Lit("c", true)
	nc := fac.Lit("c", false)

	formulas = append(formulas, cornerCasePBC3(fac, comp, []Literal{}, []int{})...)

	formulas = append(formulas, cornerCasePBC3(fac, comp, []Literal{a}, []int{-1})...)
	formulas = append(formulas, cornerCasePBC3(fac, comp, []Literal{a}, []int{0})...)
	formulas = append(formulas, cornerCasePBC3(fac, comp, []Literal{a}, []int{1})...)

	formulas = append(formulas, cornerCasePBC3(fac, comp, []Literal{na}, []int{-1})...)
	formulas = append(formulas, cornerCasePBC3(fac, comp, []Literal{na}, []int{0})...)
	formulas = append(formulas, cornerCasePBC3(fac, comp, []Literal{na}, []int{1})...)

	formulas = append(formulas, cornerCasePBC3(fac, comp, []Literal{a, b}, []int{-1, -1})...)
	formulas = append(formulas, cornerCasePBC3(fac, comp, []Literal{a, b}, []int{0, 0})...)
	formulas = append(formulas, cornerCasePBC3(fac, comp, []Literal{a, b}, []int{1, 1})...)
	formulas = append(formulas, cornerCasePBC3(fac, comp, []Literal{a, b}, []int{1, -1})...)

	formulas = append(formulas, cornerCasePBC3(fac, comp, []Literal{a, nb}, []int{-1, -1})...)
	formulas = append(formulas, cornerCasePBC3(fac, comp, []Literal{a, nb}, []int{0, 0})...)
	formulas = append(formulas, cornerCasePBC3(fac, comp, []Literal{a, nb}, []int{1, 1})...)
	formulas = append(formulas, cornerCasePBC3(fac, comp, []Literal{a, nb}, []int{1, -1})...)

	formulas = append(formulas, cornerCasePBC3(fac, comp, []Literal{a, na}, []int{-1, -1})...)
	formulas = append(formulas, cornerCasePBC3(fac, comp, []Literal{a, na}, []int{0, 0})...)
	formulas = append(formulas, cornerCasePBC3(fac, comp, []Literal{a, na}, []int{1, 1})...)
	formulas = append(formulas, cornerCasePBC3(fac, comp, []Literal{a, na}, []int{1, -1})...)

	formulas = append(formulas, cornerCasePBC3(fac, comp, []Literal{a, b, c}, []int{-1, -1, -1})...)
	formulas = append(formulas, cornerCasePBC3(fac, comp, []Literal{a, b, c}, []int{0, 0, 0})...)
	formulas = append(formulas, cornerCasePBC3(fac, comp, []Literal{a, b, c}, []int{1, 1, 1})...)
	formulas = append(formulas, cornerCasePBC3(fac, comp, []Literal{a, b, c}, []int{-1, 1, -1})...)
	formulas = append(formulas, cornerCasePBC3(fac, comp, []Literal{na, nb, nc}, []int{-1, 1, -1})...)

	return formulas
}

func cornerCasePBC3(fac Factory, comp CSort, lits []Literal, coeffs []int) []Formula {
	formulas := make([]Formula, 0, len(lits)*7)
	for _, rhs := range []int{-1, 0, 1, -3, -4, 3, 4} {
		formulas = append(formulas, fac.PBC(comp, rhs, lits, coeffs))
	}
	return formulas
}
