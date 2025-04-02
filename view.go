package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var leftPaneStyle = lipgloss.NewStyle().Width(40)

func (m Model) View() string {
	left := leftPaneStyle.Render(m.callers.View())
	right := m.paths.View()
	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}

func renderPaths(paths [][]string, target string) string {
	var b strings.Builder

	cyan := "\x1b[1;36m"   // Top-level caller
	yellow := "\x1b[1;33m" // Target function
	reset := "\x1b[0m"

	for i, path := range paths {
		fmt.Fprintf(&b, "Path #%d:\n", i+1)
		for depth, name := range path {
			prefix := strings.Repeat("  ", depth)
			if strings.EqualFold(name, target) {
				fmt.Fprintf(&b, "%s↳ %s%s%s \n", prefix, yellow, name, reset)
			} else if depth == 0 {
				fmt.Fprintf(&b, "%s↳ %s%s%s \n", prefix, cyan, name, reset)
			} else {
				fmt.Fprintf(&b, "%s↳ %s\n", prefix, name)
			}
		}
		b.WriteString("\n")
	}
	return b.String()
}
