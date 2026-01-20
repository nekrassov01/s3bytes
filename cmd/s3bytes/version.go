package main

import "fmt"

// Version is the current version of llcm.
const version = "0.0.17"

// Revision is the git revision of llcm.
var revision = ""

// version returns the version and revision of llcm.
func getVersion() string {
	if revision == "" {
		return version
	}
	return fmt.Sprintf("%s (revision: %s)", version, revision)
}
