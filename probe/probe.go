package probe

import (
	"log"
	"os"
)

// -----------------------------------------------------------------------------

// HostnamePrefix prefixes the given probe name with the hostname.
//
// This is meant to be used during software initialization: it simply crashes
// if the hostname cannot be found.
func HostnamePrefix(name string) string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	return hostname + "." + name
}
