package encoding

import (
	"testing"

	"github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/parser"
	"github.com/stretchr/testify/assert"
)

func testContainsPBC(t *testing.T) {
	assert := assert.New(t)
	fac := formula.NewFactory()
	p := parser.New(fac)

	assert.False(ContainsPBC(fac, p.ParseUnsafe("$false")))
	assert.False(ContainsPBC(fac, p.ParseUnsafe("$true")))
	assert.False(ContainsPBC(fac, p.ParseUnsafe("a")))
	assert.False(ContainsPBC(fac, p.ParseUnsafe("~a")))
	assert.False(ContainsPBC(fac, p.ParseUnsafe("~(a|b)")))
	assert.True(ContainsPBC(fac, p.ParseUnsafe("~(a | (a + b = 3))")))
	assert.True(ContainsPBC(fac, p.ParseUnsafe("~(a & (a + b = 3))")))
	assert.False(ContainsPBC(fac, p.ParseUnsafe("a => b")))
	assert.False(ContainsPBC(fac, p.ParseUnsafe("a <=> b")))
	assert.False(ContainsPBC(fac, p.ParseUnsafe("a => (b | c & ~(e | d))")))
	assert.False(ContainsPBC(fac, p.ParseUnsafe("a <=> (b | c & ~(e | d))")))
	assert.True(ContainsPBC(fac, p.ParseUnsafe("a => (3*a + ~b <= 4)")))
	assert.True(ContainsPBC(fac, p.ParseUnsafe("(3*a + ~b <= 4) <=> b")))
	assert.True(ContainsPBC(fac, p.ParseUnsafe("a => (b | c & (3*a + ~b <= 4) & ~(e | d))")))
	assert.True(ContainsPBC(fac, p.ParseUnsafe("a <=> (b | c & ~(e | (3*a + ~b <= 4) | d))")))
	assert.True(ContainsPBC(fac, p.ParseUnsafe("3*a + ~b <= 4")))
}
