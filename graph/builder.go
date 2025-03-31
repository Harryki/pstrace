package graph

import (
	"regexp"
	"strings"
)

type Builder struct {
	userFuncs map[string]struct{}
}

func NewBuilder(funcNames []string) *Builder {
	userFuncs := make(map[string]struct{})
	for _, fn := range funcNames {
		userFuncs[strings.ToLower(fn)] = struct{}{}
	}
	return &Builder{userFuncs}
}

var callRegex = regexp.MustCompile(`(?i)[A-Za-z_][A-Za-z0-9_-]*`)

func appendIfMissing(slice []string, val string) []string {
	for _, s := range slice {
		if strings.EqualFold(s, val) {
			return slice
		}
	}
	return append(slice, val)
}

func (b *Builder) Build(funcBodies map[string][]string) map[string][]string {
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
				if _, ok := b.userFuncs[strings.ToLower(callee)]; ok && !strings.EqualFold(caller, callee) {
					callGraph[callee] = appendIfMissing(callGraph[callee], caller)
				}
			}
		}
	}
	return callGraph
}
