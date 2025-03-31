package pathfinder

import "strings"

type ReverseAllPathsFinder struct{}

func NewReverseAllPathsFinder() *ReverseAllPathsFinder {
	return &ReverseAllPathsFinder{}
}

func (p *ReverseAllPathsFinder) FindPaths(graph map[string][]string, target string) [][]string {
	var results [][]string

	var dfs func(path []string, current string, visited map[string]bool)
	dfs = func(path []string, current string, visited map[string]bool) {
		if visited[strings.ToLower(current)] {
			return
		}
		visited[strings.ToLower(current)] = true
		defer delete(visited, strings.ToLower(current))

		path = append([]string{current}, path...)
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
