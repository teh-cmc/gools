package nanshal

import (
	"encoding/json"
	"math"
	"reflect"
	"strings"
)

// -----------------------------------------------------------------------------

var (
	_floatErrs = [...]string{"NaN", "+Inf", "-Inf"}
	_maxFloats = map[reflect.Kind]float64{
		reflect.Float32: math.MaxFloat32,
		reflect.Float64: math.MaxFloat64,
	}
)

// MarshalInterface marshals an interface value even if it embeds NaN, +Inf or
// -Inf float values.
//
//  - NaN values will be converted to 0.0
//  - +Inf values will be converted to _MAX_FLOAT_
//  - -Inf values will be converted to -_MAX_FLOAT_
// _MAX_FLOAT_ will be either 32 or 64 bits, depending on the struct-field that
// contains it.
//
// If the optional `maxFloat` parameter is passed, this value will be used as
// the _MAX_FLOAT_ value.
func MarshalInterface(v interface{}, maxFloat ...float64) ([]byte, error) {
	return MarshalValue(reflect.ValueOf(v), maxFloat...)
}

// MarshalValue marshals a reflect.Value even if it embeds NaN, +Inf or
// -Inf float values.
//
//  - NaN values will be converted to 0.0
//  - +Inf values will be converted to _MAX_FLOAT_
//  - -Inf values will be converted to -_MAX_FLOAT_
// _MAX_FLOAT_ will be either 32 or 64 bits, depending on the struct-field that
// contains it.
//
// If the optional `maxFloat` parameter is passed, this value will be used as
// the _MAX_FLOAT_ value.
func MarshalValue(v reflect.Value, maxFloat ...float64) (b []byte, err error) {
	vi := v.Interface()
	b, err = json.Marshal(vi)
	if err != nil {
		for _, e := range _floatErrs {
			if strings.Contains(err.Error(), e) {
				return MarshalValue(UnNaNifiyFloats(v, maxFloat...), maxFloat...)
			}
		}
		return nil, err
	}
	return b, nil
}

// -----------------------------------------------------------------------------

// UnNaNifiyFloats recursively walks through `v` and annihilates any NaN, +Inf
// or -Inf values it might find along the way.
//
//  - NaN values will be converted to 0.0
//  - +Inf values will be converted to _MAX_FLOAT_
//  - -Inf values will be converted to -_MAX_FLOAT_
// _MAX_FLOAT_ will be either 32 or 64 bits, depending on the struct-field that
// contains it.
//
// If the optional `maxFloat` parameter is passed, this value will be used as
// the _MAX_FLOAT_ value.
func UnNaNifiyFloats(v reflect.Value, maxFloat ...float64) reflect.Value {
	return unNaNifiyFloats(v, true, maxFloat...)
}
func unNaNifiyFloats(
	v reflect.Value, root bool, maxFloat ...float64,
) reflect.Value {
	// let's not adventure ourselves past the original package's boundaries
	if !root && v.Type().PkgPath() != "" {
		return v
	}
	switch v.Kind() {
	case reflect.Float32, reflect.Float64:
		vv := v.Float()
		if v.CanSet() {
			if math.IsNaN(vv) {
				v.SetFloat(0.0)
			}
			if math.IsInf(vv, -1) {
				if len(maxFloat) > 0 {
					v.SetFloat(-maxFloat[0])
				} else {
					v.SetFloat(-_maxFloats[v.Kind()])
				}
			}
			if math.IsInf(vv, 1) {
				if len(maxFloat) > 0 {
					v.SetFloat(maxFloat[0])
				} else {
					v.SetFloat(_maxFloats[v.Kind()])
				}
			}
		}
	case reflect.Ptr:
		return unNaNifiyFloats(v.Elem(), false, maxFloat...)
	case reflect.Map:
		for _, k := range v.MapKeys() {
			unNaNifiyFloats(v.MapIndex(k), false, maxFloat...)
		}
	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			unNaNifiyFloats(v.Index(i), false, maxFloat...)
		}
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			unNaNifiyFloats(v.Index(i), false, maxFloat...)
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			unNaNifiyFloats(v.Field(i), false, maxFloat...)
		}
	}
	return v
}
