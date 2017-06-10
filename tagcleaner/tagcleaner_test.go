package tagcleaner

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
)

// -----------------------------------------------------------------------------

func TestTagcleaner_Clean(t *testing.T) {
	type someStruct struct {
		A time.Time   "tag-A"
		B []time.Time "tag-B"
		C map[string]*struct {
			A string "tag-C-A"
			B string "tag-C-B"
		} "tag-C"
		D map[int32]time.Duration "tag-D"
		E struct {
			A int64 "tag-E-A"
			B int32 "tag-E-B"
		} "tag-E"
		F *struct {
			A *struct {
				A int64 "tag-F-A-A"
				B int32 "tag-F-A-B"
			} "tag-F-A"
		} "tag-F"
		G struct {
			A string "tag-G-A"
			B string "tag-G-B"
		} "tag-G"
		H [3]*struct {
			A []struct {
				A int "tag-H-A-A"
			}
			B string "tag-H-B"
			C string "tag-H-C"
		} "tag-G"
	}

	raw := fmt.Sprintf("%s", reflect.TypeOf(someStruct{}))
	if !strings.Contains(raw, "tag") {
		t.Error("raw dump should contain tags")
	}
	cleansed := fmt.Sprintf("%s", Clean(reflect.TypeOf(someStruct{})))
	if strings.Contains(cleansed, "tag") {
		t.Error("cleansed dump should NOT contain tags")
	}
	if !strings.Contains(cleansed, "[]") {
		t.Error("cleansed dump should still contain slice markers")
	}
	if !strings.Contains(cleansed, "3]") {
		t.Error("cleansed dump should still contain array markers")
	}
	if !strings.Contains(cleansed, "*") {
		t.Error("cleansed dump should still contain pointer markers")
	}
}
