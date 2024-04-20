package transformation

import (
	"testing"

	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/parser"
	"github.com/stretchr/testify/assert"
)

func TestAnonymizerWithoutPrefix(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)

	a1 := Anonymize(fac, p.ParseUnsafe("a & b & (a | b | c)"))
	assert.Equal(p.ParseUnsafe("v0 & v1 & (v0 | v1 | v2)"), a1)

	anonymizer := NewAnonymizer(fac)
	a1 = anonymizer.Anonymize(p.ParseUnsafe("a & b & (a | b | c)"))
	assert.Equal(p.ParseUnsafe("v0 & v1 & (v0 | v1 | v2)"), a1)
	a1 = anonymizer.Anonymize(p.ParseUnsafe("a & ~c"))
	assert.Equal(p.ParseUnsafe("v0 & ~v2"), a1)
}

func TestAnonymizerWithPrefix(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)

	anonymizer := NewAnonymizer(fac, "x")
	a1 := anonymizer.Anonymize(p.ParseUnsafe("a & b & (a | b | c)"))
	assert.Equal(p.ParseUnsafe("x0 & x1 & (x0 | x1 | x2)"), a1)
	a1 = anonymizer.Anonymize(p.ParseUnsafe("a & ~c"))
	assert.Equal(p.ParseUnsafe("x0 & ~x2"), a1)
}
