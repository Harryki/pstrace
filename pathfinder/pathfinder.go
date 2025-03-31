package pathfinder

type PathFinder interface {
	FindPaths(graph map[string][]string, target string) [][]string
}
