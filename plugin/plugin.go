package main

import (
	"github.com/dontpanicw/log-linter"
	"golang.org/x/tools/go/analysis"
)

// AnalyzerPlugin is required for golangci-lint plugin
type AnalyzerPlugin struct{}

// GetAnalyzers returns the list of analyzers
func (*AnalyzerPlugin) GetAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		loglinter.Analyzer,
	}
}

// New creates a new instance of the plugin
var New AnalyzerPlugin
