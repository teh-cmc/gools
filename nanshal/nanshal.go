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
func MarshalInterface(v interface{}) ([]byte, error) {
	return MarshalValue(reflect.ValueOf(v))
}

// MarshalValue marshals a reflect.Value even if it embeds NaN, +Inf or
// -Inf float values.
//
//  - NaN values will be converted to 0.0
//  - +Inf values will be converted to _MAX_FLOAT_
//  - -Inf values will be converted to -_MAX_FLOAT_
func MarshalValue(v reflect.Value) (b []byte, err error) {
	vi := v.Interface()
	b, err = json.Marshal(vi)
	if err != nil {
		for _, e := range _floatErrs {
			if strings.Contains(err.Error(), e) {
				return MarshalValue(UnNaNifiyFloats(v))
			}
		}
		return nil, err
	}
	return b, nil
}

// -----------------------------------------------------------------------------

// UnNaNifiyFloats recursively walks through `v` and annihilates any NaN, +Inf
// or -Inf values it might find along the way.
func UnNaNifiyFloats(v reflect.Value) reflect.Value {
	return unNaNifiyFloats(v, true)
}
func unNaNifiyFloats(v reflect.Value, root bool) reflect.Value {
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
				v.SetFloat(-_maxFloats[v.Kind()])
			}
			if math.IsInf(vv, 1) {
				v.SetFloat(_maxFloats[v.Kind()])
			}
		}
	case reflect.Ptr:
		return unNaNifiyFloats(v.Elem(), false)
	case reflect.Map:
		for _, k := range v.MapKeys() {
			unNaNifiyFloats(v.MapIndex(k), false)
		}
	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			unNaNifiyFloats(v.Index(i), false)
		}
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			unNaNifiyFloats(v.Index(i), false)
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			unNaNifiyFloats(v.Field(i), false)
		}
	}
	return v
}
