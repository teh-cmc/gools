package probe

import (
	"os"
	"testing"
)

// -----------------------------------------------------------------------------

func TestProbe_HostnamePrefix(t *testing.T) {
	probe := "queries.processed"
	hostname, err := os.Hostname()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	expected := hostname + "." + probe
	prefixed := HostnamePrefix(probe)
	if prefixed != expected {
		t.Errorf("expected %v, got %v", expected, prefixed)
	}
}
