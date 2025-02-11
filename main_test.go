package s3bytes

import (
	"testing"
)

// TestMain is the entry point of the test.
func TestMain(m *testing.M) {
	var (
		originalMaxQueries    = setMaxQueries(2)
		originalMaxChartItems = setMaxChartItems(3)
	)
	defer func() {
		setMaxQueries(originalMaxQueries)
		setMaxChartItems(originalMaxChartItems)
	}()
	m.Run()
}

func setMaxQueries(n int) (original int) {
	original = MaxQueries
	MaxQueries = n
	return original
}

func setMaxChartItems(n int) (original int) {
	original = MaxChartItems
	MaxChartItems = n
	return original
}
