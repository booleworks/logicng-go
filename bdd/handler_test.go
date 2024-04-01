package bdd

import (
	"testing"
	"time"

	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/handler"
	"github.com/booleworks/logicng-go/normalform"
	"github.com/booleworks/logicng-go/sat"
	"github.com/stretchr/testify/assert"
)

func TestBDDTimeoutHandlerSmall(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	queens := sat.GenerateNQueens(fac, 4)
	kernel := NewKernel(fac, int32(f.Variables(fac, queens).Size()), 10000, 10000)
	duration, _ := time.ParseDuration("2s")
	handler := HandlerWithTimeout(*handler.NewTimeoutWithDuration(duration))
	bdd, ok := BuildWithKernelAndHandler(fac, queens, kernel, handler)
	assert.True(ok)
	assert.False(handler.Aborted())
	assert.Greater(bdd.Index, int32(0))
}

func TestBDDTimeoutHandlerLarge(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	queens := normalform.CNF(fac, sat.GenerateNQueens(fac, 15))
	kernel := NewKernel(fac, int32(f.Variables(fac, queens).Size()), 100000, 100000)
	duration, _ := time.ParseDuration("500ms")
	handler := HandlerWithTimeout(*handler.NewTimeoutWithDuration(duration))
	bdd, ok := BuildWithKernelAndHandler(fac, queens, kernel, handler)
	assert.False(ok)
	assert.True(handler.Aborted())
	assert.Nil(bdd)
}

func TestBDDNodesHandlerSmall(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	queens := sat.GenerateNQueens(fac, 4)
	kernel := NewKernel(fac, int32(f.Variables(fac, queens).Size()), 10000, 10000)
	handler := HandlerWithNodes(1000)
	bdd, ok := BuildWithKernelAndHandler(fac, queens, kernel, handler)
	assert.True(ok)
	assert.False(handler.Aborted())
	assert.Greater(bdd.Index, int32(0))
}

func TestBDDNodesHandlerLarge(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	queens := sat.GenerateNQueens(fac, 25)
	kernel := NewKernel(fac, int32(f.Variables(fac, queens).Size()), 10000, 10000)
	handler := HandlerWithNodes(50)
	bdd, ok := BuildWithKernelAndHandler(fac, queens, kernel, handler)
	assert.False(ok)
	assert.True(handler.Aborted())
	assert.Nil(bdd)
}
