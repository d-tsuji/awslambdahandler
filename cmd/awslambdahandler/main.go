package main

import (
	"github.com/d-tsuji/awslambdahandler"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(awslambdahandler.Analyzer) }
