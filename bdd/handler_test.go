package bdd

import (
	"testing"
	"time"

	"github.com/booleworks/logicng-go/event"
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
	hdl := handler.NewTimeoutWithDuration(duration)
	bdd, state := CompileWithKernelAndHandler(fac, queens, kernel, hdl)
	assert.True(state.Success)
	assert.Equal(event.Nothing, state.CancelCause)
	assert.Greater(bdd.Index, int32(0))
}

func TestBDDTimeoutHandlerLarge(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	queens := normalform.CNF(fac, sat.GenerateNQueens(fac, 15))
	kernel := NewKernel(fac, int32(f.Variables(fac, queens).Size()), 100000, 100000)
	duration, _ := time.ParseDuration("500ms")
	hdl := handler.NewTimeoutWithDuration(duration)
	bdd, state := CompileWithKernelAndHandler(fac, queens, kernel, hdl)
	assert.False(state.Success)
	assert.Equal(event.BddNewRefAdded, state.CancelCause)
	assert.Nil(bdd)
}

func TestBDDNodesHandlerSmall(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	queens := sat.GenerateNQueens(fac, 4)
	kernel := NewKernel(fac, int32(f.Variables(fac, queens).Size()), 10000, 10000)
	hdl := HandlerWithNodes(1000)
	bdd, state := CompileWithKernelAndHandler(fac, queens, kernel, hdl)
	assert.True(state.Success)
	assert.Equal(event.Nothing, state.CancelCause)
	assert.Greater(bdd.Index, int32(0))
}

func TestBDDNodesHandlerLarge(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	queens := sat.GenerateNQueens(fac, 25)
	kernel := NewKernel(fac, int32(f.Variables(fac, queens).Size()), 10000, 10000)
	hdl := HandlerWithNodes(50)
	bdd, state := CompileWithKernelAndHandler(fac, queens, kernel, hdl)
	assert.False(state.Success)
	assert.Equal(event.BddNewRefAdded, state.CancelCause)
	assert.Nil(bdd)
}
