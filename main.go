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
			// if line starts with #, ignore
			// if line starts with <#, set ignoring to true.
			// if line starts with #>, ig nore the line, check if ingoring is true, if so, set it to false

			if strings.Contains(line, "{") {
				counter += strings.Count(line, "{")
			}
			if strings.Contains(line, "}") {
				counter -= strings.Count(line, "}")
			}

			if counter > 0 {
				buffer = append(buffer, line)
			} else if counter == 0 && buffer != nil {
				// function body is complete
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
	for name, body := range funcBodies {
		for _, line := range body {
			for _, match := range callRegex.FindAllStringIndex(line, -1) {
				start := match[0]
				end := match[1]
				word := line[start:end]

				// Check character before the match
				if start > 0 {
					prev := line[start-1]
					if prev == '$' || prev == '-' || prev == '.' {
						continue // not a function call
					}
				}

				//  if word is found in a funcNames and is not the same as the current function
				if _, ok := userFuncs[strings.ToLower(word)]; ok && !strings.EqualFold(name, word) {
					callGraph[name] = appendIfMissing(callGraph[name], word)
				}
			}
		}

	}

	// Print call graph
	for caller, callees := range callGraph {
		for _, callee := range callees {
			fmt.Printf("%s -> %s\n", caller, callee)
		}
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
