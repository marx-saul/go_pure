package main

import (
	"github.com/marx-saul/go_pure"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(go_pure.Analyzer) }

