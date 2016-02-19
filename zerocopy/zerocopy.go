package zerocopy

import (
	"reflect"
	"unsafe"
)

// -----------------------------------------------------------------------------

// StringToByteSlice converts a string to a []byte without any copy.
//
// This is obviously particularly unsafe and should be used at your own risk.
//
// NOTE: do not ever modify the returned byte slice.
// NOTE: do not ever use the returned byte slice once the original string went
// out of scope.
func StringToByteSlice(s string) (b []byte) {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bh.Len, bh.Cap, bh.Data = sh.Len, sh.Len, sh.Data
	return
}
