package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	funcDeclRegex = regexp.MustCompile(`(?i)^function\s+([A-Za-z_][A-Za-z0-9_-]*)\s*(\(\))?\s*\{?`)
	callRegex     = regexp.MustCompile(`(?i)[A-Za-z_][A-Za-z0-9_-]*`)
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
	lines := strings.Split(string(content), "\n")

	// Step 1: Find function declarations and their line ranges
	funcBodies := map[string][]string{}
	funcNames := []string{}

	var currentFunc string

	var collecting bool
	var buffer []string
	var counter int
	var ignoring bool

	for _, line := range lines {
		lineTrim := strings.TrimSpace(line)

		// If line matches a function declaration pattern, we will collect function name and its full body
		if matches := funcDeclRegex.FindStringSubmatch(lineTrim); matches != nil {
			currentFunc = matches[1]

			funcNames = append(funcNames, currentFunc)
			// fmt.Printf("Parsed function: %s\n", currentFunc)
			collecting = true
			buffer = nil
		}

		if collecting {
			if strings.Contains(lineTrim, "<#") && strings.Contains(lineTrim, "#>") {
				continue // one-line block comment
			}
			if strings.HasPrefix(lineTrim, "<#") {
				ignoring = true
				continue
			}
			if strings.Contains(lineTrim, "#>") {
				ignoring = false
				continue
			}
			if ignoring || strings.HasPrefix(lineTrim, "#") {
				continue
			}
			if idx := strings.Index(lineTrim, "#"); idx != -1 {
				lineTrim = strings.TrimSpace(lineTrim[:idx])
			}
			if lineTrim == "" {
				continue
			}

			if strings.Contains(line, "{") {
				counter += strings.Count(line, "{")
			}
			if strings.Contains(line, "}") {
				counter -= strings.Count(line, "}")
			}

			if counter > 0 {
				buffer = append(buffer, lineTrim)
			}

			if counter == 0 && buffer != nil {
				funcBodies[currentFunc] = append([]string{}, buffer...)
				buffer = nil
				collecting = false
			}
		}
	}
	if counter != 0 {
		panic("Counter is negative when it should be 0!")
	}

	// Step 2: Build call graph
	userFuncs := map[string]struct{}{}
	for _, fn := range funcNames {
		userFuncs[strings.ToLower(fn)] = struct{}{}
	}

	callGraph := map[string][]string{}
	for caller, body := range funcBodies {
		for _, line := range body {
			for _, match := range callRegex.FindAllStringIndex(line, -1) {
				start := match[0]
				end := match[1]
				callee := line[start:end]

				// Check character before the match
				if start > 0 {
					prev := line[start-1]
					if prev == '$' || prev == '-' || prev == '.' {
						continue // not a function call
					}
				}

				//  if word is found in a funcNames and is not the same as the current function
				if _, ok := userFuncs[strings.ToLower(callee)]; ok && !strings.EqualFold(caller, callee) {
					callGraph[callee] = appendIfMissing(callGraph[callee], caller)
				}
			}
		}

	}

	// Print call graph
	// for caller, callees := range callGraph {
	// 	for _, callee := range callees {
	// 		fmt.Printf("%s -> %s\n", caller, callee)
	// 	}
	// }

	// Check for --paths-to option
	if len(os.Args) >= 4 && os.Args[2] == "--paths-to" {
		target := os.Args[3]
		paths := findPathsToTargetReversed(callGraph, target)
		if len(paths) == 0 {
			fmt.Printf("No paths found to %s\n", target)
			return
		}
		printPathsIndented(paths, target)
		return
	}
}

func appendIfMissing(slice []string, val string) []string {
	for _, s := range slice {
		if strings.EqualFold(s, val) {
			return slice
		}
	}
	return append(slice, val)
}

func findPathsToTargetReversed(graph map[string][]string, target string) [][]string {
	var results [][]string

	var dfs func(path []string, current string, visited map[string]bool)
	dfs = func(path []string, current string, visited map[string]bool) {
		if visited[strings.ToLower(current)] {
			return
		}
		visited[strings.ToLower(current)] = true
		defer delete(visited, strings.ToLower(current))

		path = append([]string{current}, path...) // prepend
		callers := graph[current]
		if len(callers) == 0 {
			results = append(results, append([]string{}, path...))
			return
		}
		for _, caller := range callers {
			dfs(path, caller, visited)
		}
	}

	dfs([]string{}, target, map[string]bool{})
	return results
}

func printPathsIndented(paths [][]string, target string) {
	yellow := "\x1b[1;33m"
	reset := "\x1b[0m"

	for i, path := range paths {
		fmt.Printf("Path #%d:\n", i+1)
		for depth, name := range path {
			prefix := strings.Repeat("  ", depth)
			if strings.EqualFold(name, target) {
				fmt.Printf("%s↳ %s%s%s  ← target\n", prefix, yellow, name, reset)
			} else {
				fmt.Printf("%s↳ %s\n", prefix, name)
			}
		}
		fmt.Println()
	}
}
