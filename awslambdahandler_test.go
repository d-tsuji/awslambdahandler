package awslambdahandler_test

import (
	"github.com/d-tsuji/awslambdahandler"
	"github.com/gostaticanalysis/testutil"
	"golang.org/x/tools/go/analysis/analysistest"

	"testing"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := testutil.WithModules(t, analysistest.TestData(), nil)
	analysistest.Run(t, testdata, awslambdahandler.Analyzer, "a")
}
