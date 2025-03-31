package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/harryki/pstrace/graph"
	"github.com/harryki/pstrace/parser"
	"github.com/harryki/pstrace/pathfinder"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: pstrace file.psm1")
		return
	}

	content, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	parser := parser.NewRegexParser()
	funcBodies, funcNames, err := parser.ParseFunctions(string(content))
	if err != nil {
		panic(err)
	}

	builder := graph.NewBuilder(funcNames)
	callGraph := builder.Build(funcBodies)

	// e.g. --paths-to "Get-MocConfig" will printout all invokation hierarchy that calls Get-MocConfig
	if len(os.Args) >= 4 && os.Args[2] == "--paths-to" {
		target := os.Args[3]
		pf := pathfinder.NewReverseAllPathsFinder()
		paths := pf.FindPaths(callGraph, target)

		m := NewCallGraphTUIModel(paths, target) // you'll define this in model.go
		if _, err := tea.NewProgram(m).Run(); err != nil {
			os.Exit(1)
		}

		return
	}
}
