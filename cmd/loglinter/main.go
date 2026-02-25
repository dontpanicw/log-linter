package main

import (
	"github.com/dontpanicw/log-linter"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(loglinter.Analyzer)
}
