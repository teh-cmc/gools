package nanshal

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// -----------------------------------------------------------------------------

type Floats struct {
	X  float64
	Y  int32
	ZZ *float64
	FF []float64
	UU struct {
		KK float64
		II []float64
	}
	SS string
}

func TestNanshal_MarshalValue(t *testing.T) {
	zz := math.NaN()
	before := Floats{
		X:  math.NaN(),
		Y:  42,
		ZZ: &zz,
		FF: []float64{math.NaN(), math.NaN(), math.NaN()},
		UU: struct {
			KK float64
			II []float64
		}{
			KK: math.NaN(),
			II: []float64{math.NaN()},
		},
		SS: "coucou",
	}
	expected := Floats{
		X:  0.0,
		Y:  42,
		ZZ: &zz,
		FF: []float64{0.0, 0.0, 0.0},
		UU: struct {
			KK float64
			II []float64
		}{
			KK: 0.0,
			II: []float64{0.0},
		},
		SS: "coucou",
	}

	afterB, err := MarshalInterface(&before)
	assert.NoError(t, err)
	expectedB, err := json.Marshal(&expected)
	assert.NoError(t, err)
	assert.Equal(t, expectedB, afterB)
}
