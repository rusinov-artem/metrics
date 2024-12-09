package main

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/analysis/code"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
)

func main() {
	alist := []*analysis.Analyzer{}

	// SA checkers from staticcheck.io
	for _, v := range staticcheck.Analyzers {
		alist = append(alist, v.Analyzer)
	}

	// S checkers from staticcheck.io
	for _, v := range simple.Analyzers {
		alist = append(alist, v.Analyzer)
	}

	// some random stuff
	alist = append(alist, printf.Analyzer)
	alist = append(alist, structtag.Analyzer)
	alist = append(alist, myanalyzer)

	multichecker.Main(
		alist[len(alist)-1:]...,
	)
}

var myanalyzer = &analysis.Analyzer{
	Name: "myanalyzer",
	Doc:  "check my stuff",
	Run: func(pass *analysis.Pass) (interface{}, error) {
		var fmain ast.Node
		for _, file := range pass.Files {
			ast.Inspect(file, func(node ast.Node) bool {
				switch node := node.(type) {
				case *ast.FuncDecl:
					if node.Name.Name != "main" {
						return false
					}
					fmain = node
					return true
				case *ast.CallExpr:
					if file.Name.Name != "main" {
						return false
					}
					if code.IsCallTo(pass, node, "os.Exit") && fmain != nil {
						pass.Reportf(node.Pos(), "os.Exist should not be called directly from main")
						return false
					}
				}
				return true
			})
		}
		return nil, nil
	},
}
