package normalform

import (
	"testing"
	"time"

	"github.com/booleworks/logicng-go/event"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/handler"
	"github.com/booleworks/logicng-go/randomizer"
	"github.com/stretchr/testify/assert"
)

func TestCNFFactorizationTimeoutHandlerSmall(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	formula := randomizer.New(fac).Formula(2)
	duration, _ := time.ParseDuration("2s")
	hdl := handler.NewTimeoutWithDuration(duration)
	cnf, state := FactorizedCNFWithHandler(fac, formula, hdl)
	assert.True(state.Success)
	assert.Equal(event.Nothing, state.CancelCause)
	assert.NotEqual(fac.Falsum(), cnf)
}

func TestCNFFactorizationTimeoutHandlerLarge(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	formula := randomizer.New(fac).Formula(7)
	duration, _ := time.ParseDuration("5ms")
	hdl := handler.NewTimeoutWithDuration(duration)
	cnf, state := FactorizedCNFWithHandler(fac, formula, hdl)
	assert.False(state.Success)
	assert.NotEqual(event.Nothing, state.CancelCause)
	assert.Equal(fac.Falsum(), cnf)
}

func TestDNFFactorizationTimeoutHandlerSmall(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	formula := randomizer.New(fac).Formula(2)
	duration, _ := time.ParseDuration("2s")
	hdl := handler.NewTimeoutWithDuration(duration)
	dnf, state := FactorizedDNFWithHandler(fac, formula, hdl)
	assert.True(state.Success)
	assert.Equal(event.Nothing, state.CancelCause)
	assert.NotEqual(fac.Falsum(), dnf)
}

func TestDNFFactorizationTimeoutHandlerLarge(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	formula := randomizer.New(fac).Formula(7)
	duration, _ := time.ParseDuration("5ms")
	hdl := handler.NewTimeoutWithDuration(duration)
	dnf, state := FactorizedDNFWithHandler(fac, formula, hdl)
	assert.False(state.Success)
	assert.NotEqual(event.Nothing, state.CancelCause)
	assert.Equal(fac.Falsum(), dnf)
}
