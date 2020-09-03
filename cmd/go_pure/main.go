package main

import (
	"go_pure"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(go_pure.Analyzer) }

