package main

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
)

type Model struct {
	callers    list.Model
	paths      viewport.Model
	callerData map[string][][]string // caller -> list of paths to target
	target     string
}

func NewCallGraphTUIModel(paths [][]string, target string) Model {
	callers := []list.Item{}
	seen := map[string]bool{}
	mockPaths := map[string][][]string{}

	for _, path := range paths {
		if len(path) == 0 {
			continue
		}
		first := path[0]
		if !seen[first] {
			callers = append(callers, listItem(first))
			seen[first] = true
			mockPaths[first] = [][]string{path}
		} else {
			mockPaths[first] = append(mockPaths[first], path)
		}
	}

	// height := len(callers) + 2 // +2 for title + padding
	listModel := list.New(callers, list.NewDefaultDelegate(), 25, 50)
	listModel.Title = "Top-Level Callers"
	vp := viewport.New(60, 40)
	vp.SetContent(renderPaths(mockPaths[callers[0].FilterValue()], target))

	return Model{
		callers:    listModel,
		paths:      vp,
		callerData: mockPaths,
		target:     target,
	}
}

type listItem string

func (i listItem) Title() string       { return string(i) }
func (i listItem) Description() string { return "" }
func (i listItem) FilterValue() string { return string(i) }
