package env

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/fatih/camelcase"
)

// -----------------------------------------------------------------------------

// String takes any structure as input and returns its string-representation
// following the classic ENV format as output.
//
// if anything should go wrong, it returns an empty string.
func String(st interface{}) string {
	cv := reflect.ValueOf(st)
	fields := reflect.TypeOf(st)
	if fields.Kind() == reflect.Ptr {
		return String(cv.Elem().Interface())
	} else if fields.Kind() == reflect.Struct {
	} else {
		return ""
	}

	fieldStrs := make([]string, fields.NumField())
	for i := 0; i < fields.NumField(); i++ {
		f := fields.Field(i)
		if len(f.Name) > 0 && strings.ToLower(f.Name[:1]) == f.Name[:1] {
			continue // private field
		}
		fNameParts := camelcase.Split(f.Name)
		for i, fnp := range fNameParts {
			fNameParts[i] = strings.ToUpper(fnp)
		}
		fName := strings.Join(fNameParts, "_")
		itf := cv.Field(i).Interface()
		if _, ok := itf.(fmt.Stringer); ok {
			fieldStrs[i] = fmt.Sprintf("%s = %s", fName, itf)
		} else {
			fieldStrs[i] = fmt.Sprintf("%s = %v", fName, itf)
		}
	}
	return strings.Join(fieldStrs, "\n")
}
