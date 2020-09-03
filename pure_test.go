package pure_test

import (
	"testing"

	"github.com/marx-saul/go-pure"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, pure.Analyzer, "a")
}

