package main

import (
	"pure"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(pure.Analyzer) }

