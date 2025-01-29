package s3bytes

import (
	"os"
	"testing"
)

// TestMain is the entry point of the test.
func TestMain(m *testing.M) {
	original := MaxQueries
	MaxQueries = 2
	code := m.Run()
	defer func() {
		MaxQueries = original
		os.Exit(code)
	}()
}
