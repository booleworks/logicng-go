package encoding

import (
	"strings"
	"testing"

	f "github.com/booleworks/logicng-go/formula"
	"github.com/stretchr/testify/assert"
)

var configs = []Config{
	{AMOEncoder: AMOPure},
	{AMOEncoder: AMOLadder},
	{AMOEncoder: AMOBinary},
	{AMOEncoder: AMOProduct},
	{AMOEncoder: AMOProduct, ProductRecursiveBound: 10},
	{AMOEncoder: AMONested},
	{AMOEncoder: AMONested, NestingGroupSize: 5},
	{AMOEncoder: AMOCommander, CommanderGroupSize: 3},
	{AMOEncoder: AMOCommander, CommanderGroupSize: 7},
	{AMOEncoder: AMOBimander, BimanderGroupSize: BimanderFixed},
	{AMOEncoder: AMOBimander, BimanderGroupSize: BimanderHalf},
	{AMOEncoder: AMOBimander, BimanderGroupSize: BimanderSqrt},
	{AMOEncoder: AMOBimander, BimanderGroupSize: BimanderFixed, BimanderFixedGroupSize: 2},
	{AMOEncoder: AMOBest},
}

func TestAMOZero(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	cc := fac.AMO()

	assert.Equal(fac.Verum(), cc)
}

func TestAMOOne(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	cc := fac.AMO(fac.Var("v0"))

	for _, config := range configs {
		cnf, err := EncodeCC(fac, cc, &config)
		assert.Nil(err)
		assert.Equal(0, len(cnf))
	}
	name, _ := fac.VarName(fac.NewCCVariable())
	assert.True(strings.HasSuffix(name, "_0"))
}
