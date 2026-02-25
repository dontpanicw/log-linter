package loglinter_test

import (
	"testing"

	"github.com/dontpanicw/log-linter"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, loglinter.Analyzer, "a")
}
