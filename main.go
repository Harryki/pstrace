package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

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

		// sort paths by first element
		sortPathsByRoot(paths)

		printPathsIndented(paths, target)
		return
	}
}

func sortPathsByRoot(paths [][]string) {
	sort.Slice(paths, func(i, j int) bool {
		// compare case-insensitively
		left := strings.ToLower(paths[i][0])
		right := strings.ToLower(paths[j][0])
		return left < right
	})
}

func printPathsIndented(paths [][]string, target string) {
	yellow := "\x1b[1;33m"
	reset := "\x1b[0m"

	for i, path := range paths {
		fmt.Printf("Path #%d:\n", i+1)
		for depth, name := range path {
			prefix := strings.Repeat("  ", depth)
			if strings.EqualFold(name, target) {
				fmt.Printf("%s↳ %s%s%s \n", prefix, yellow, name, reset)
			} else {
				fmt.Printf("%s↳ %s\n", prefix, name)
			}
		}
		fmt.Println()
	}
}
