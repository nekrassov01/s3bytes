package s3bytes

import (
	"os"
	"testing"
)

// TestMain is the entry point of the test.
func TestMain(m *testing.M) {
	originalMaxQueries := MaxQueries
	MaxQueries = 2
	originalMaxChartItems := MaxChartItems
	MaxChartItems = 3
	code := m.Run()
	defer func() {
		MaxQueries = originalMaxQueries
		MaxChartItems = originalMaxChartItems
		os.Exit(code)
	}()
}
