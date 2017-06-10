package tagcleaner

import (
	"reflect"
)

// -----------------------------------------------------------------------------

// Clean recursively walks through `t` and annihilates any struct-tag it might
// find along the way.
func Clean(t reflect.Type) reflect.Type { return clean(t, true) }

func clean(t reflect.Type, root bool) reflect.Type {
	// let's not adventure ourselves past the original package's boundaries
	if !root && t.PkgPath() != "" {
		return t
	}
	switch t.Kind() {
	case reflect.Ptr:
		return reflect.PtrTo(clean(t.Elem(), false))
	case reflect.Map:
		return reflect.MapOf(t.Key(), clean(t.Elem(), false))
	case reflect.Array:
		return reflect.ArrayOf(t.Len(), clean(t.Elem(), false))
	case reflect.Slice:
		return reflect.SliceOf(clean(t.Elem(), false))
	case reflect.Struct:
		fields := make([]reflect.StructField, t.NumField())
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			f.Tag = ""
			f.Type = clean(f.Type, false)
			fields[i] = f
		}
		return reflect.StructOf(fields)
	default:
		return t
	}
}
