package normalform

import (
	"testing"
	"time"

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
	hdl := HandlerWithTimeout(*handler.NewTimeoutWithDuration(duration))
	cnf, ok := FactorizedCNFWithHandler(fac, formula, hdl)
	assert.True(ok)
	assert.False(hdl.Aborted())
	assert.NotEqual(fac.Falsum(), cnf)
}

func TestCNFFactorizationTimeoutHandlerLarge(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	formula := randomizer.New(fac).Formula(7)
	duration, _ := time.ParseDuration("5ms")
	hdl := HandlerWithTimeout(*handler.NewTimeoutWithDuration(duration))
	cnf, ok := FactorizedCNFWithHandler(fac, formula, hdl)
	assert.False(ok)
	assert.True(hdl.Aborted())
	assert.Equal(fac.Falsum(), cnf)
}

func TestDNFFactorizationTimeoutHandlerSmall(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	formula := randomizer.New(fac).Formula(2)
	duration, _ := time.ParseDuration("2s")
	hdl := HandlerWithTimeout(*handler.NewTimeoutWithDuration(duration))
	dnf, ok := FactorizedDNFWithHandler(fac, formula, hdl)
	assert.True(ok)
	assert.False(hdl.Aborted())
	assert.NotEqual(fac.Falsum(), dnf)
}

func TestDNFFactorizationTimeoutHandlerLarge(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	formula := randomizer.New(fac).Formula(7)
	duration, _ := time.ParseDuration("5ms")
	hdl := HandlerWithTimeout(*handler.NewTimeoutWithDuration(duration))
	dnf, ok := FactorizedDNFWithHandler(fac, formula, hdl)
	assert.False(ok)
	assert.True(hdl.Aborted())
	assert.Equal(fac.Falsum(), dnf)
}
