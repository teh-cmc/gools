package tagcleaner

import (
	"reflect"
)

// -----------------------------------------------------------------------------

// Clean recursively walks through `t` and annihilates any struct-tag it might
// find along the way.
func Clean(t reflect.Type) reflect.Type {
	switch t.Kind() {
	case reflect.Ptr:
		return reflect.PtrTo(Clean(t.Elem()))
	case reflect.Map:
		return reflect.MapOf(t.Key(), Clean(t.Elem()))
	case reflect.Array:
		return reflect.ArrayOf(t.Len(), Clean(t.Elem()))
	case reflect.Slice:
		return reflect.SliceOf(Clean(t.Elem()))
	case reflect.Struct:
		fields := make([]reflect.StructField, t.NumField())
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			f.Tag = ""
			f.Type = Clean(f.Type)
			fields[i] = f
		}
		return reflect.StructOf(fields)
	default:
		return t
	}
}
