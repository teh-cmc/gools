package nanshal

import (
	"encoding/json"
	"reflect"
	"strings"
)

// -----------------------------------------------------------------------------

// MarshalInterface marshals an interface value even if it embeds NaN floats.
//
// NaN values will be converted to 0.0.
func MarshalInterface(v interface{}) ([]byte, error) {
	return MarshalValue(reflect.ValueOf(v))
}

// MarshalValue marshals a reflect.Value even if it embeds NaN floats.
//
// NaN values will be converted to 0.0.
func MarshalValue(v reflect.Value) (b []byte, err error) {
	vi := v.Interface()
	b, err = json.Marshal(vi)
	if err != nil {
		if strings.Contains(err.Error(), "NaN") {
			return MarshalValue(UnNaNifiyFloats(v))
		}
		return nil, err
	}
	return b, nil
}

// -----------------------------------------------------------------------------

// UnNaNifiyFloats recursively walks through `v` and annihilates any NaN values
// it might find along the way.
func UnNaNifiyFloats(v reflect.Value) reflect.Value {
	return unNaNifiyFloats(v, true)
}
func unNaNifiyFloats(v reflect.Value, root bool) reflect.Value {
	// let's not adventure ourselves past the original package's boundaries
	if !root && v.Type().PkgPath() != "" {
		return v
	}
	switch v.Kind() {
	case reflect.Float32:
		if v.CanSet() {
			v.Set(reflect.ValueOf(float32(0)))
		}
	case reflect.Float64:
		if v.CanSet() {
			v.Set(reflect.ValueOf(float64(0)))
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
