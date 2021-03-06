package nanshal

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// -----------------------------------------------------------------------------

type Floats struct {
	A  float32
	B  float64
	I  float32
	J  float32
	K  float64
	L  float64
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
		A:  1.0,
		B:  2.0,
		I:  float32(math.Inf(-1)),
		J:  float32(math.Inf(1)),
		K:  float64(math.Inf(-1)),
		L:  float64(math.Inf(1)),
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
		A:  1.0,
		B:  2.0,
		I:  float32(-math.MaxFloat32),
		J:  float32(math.MaxFloat32),
		K:  float64(-math.MaxFloat64),
		L:  float64(math.MaxFloat64),
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

func TestNanshal_MarshalValue_MaxFloat(t *testing.T) {
	var maxFloat float64 = 42.0
	zz := math.NaN()
	before := Floats{
		A:  1.0,
		B:  2.0,
		I:  float32(math.Inf(-1)),
		J:  float32(math.Inf(1)),
		K:  float64(math.Inf(-1)),
		L:  float64(math.Inf(1)),
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
		A:  1.0,
		B:  2.0,
		I:  float32(-maxFloat),
		J:  float32(maxFloat),
		K:  float64(-maxFloat),
		L:  float64(maxFloat),
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

	afterB, err := MarshalInterface(&before, maxFloat)
	assert.NoError(t, err)
	expectedB, err := json.Marshal(&expected)
	assert.NoError(t, err)
	assert.Equal(t, expectedB, afterB)
}
