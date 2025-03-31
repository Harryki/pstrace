package main

import (
	"sort"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
)

type Model struct {
	callers    list.Model
	paths      viewport.Model
	callerData map[string][][]string // caller -> list of paths to target
	target     string
}

type listItem string

func (i listItem) Title() string       { return string(i) }
func (i listItem) Description() string { return "" }
func (i listItem) FilterValue() string { return string(i) }

func NewCallGraphTUIModel(paths [][]string, target string) Model {
	callers := []list.Item{}
	seen := map[string]bool{}
	pathsBycaller := map[string][][]string{}

	for _, path := range paths {
		if len(path) == 0 {
			continue
		}
		first := path[0]
		if !seen[first] {
			callers = append(callers, listItem(first))
			seen[first] = true
			pathsBycaller[first] = [][]string{path}
		} else {
			pathsBycaller[first] = append(pathsBycaller[first], path)
		}
	}
	// TODO: maybe create callers first with sorted order and build pathsBycaller?
	// Convert []list.Item to []listItem
	var rawItems []listItem
	for _, item := range callers {
		rawItems = append(rawItems, item.(listItem))
	}

	// Sort by string value
	sort.Slice(rawItems, func(i, j int) bool {
		return rawItems[i] < rawItems[j] // or use strings.ToLower(...) if needed
	})

	// Convert back to []list.Item
	callers = make([]list.Item, len(rawItems))
	for i, item := range rawItems {
		callers[i] = item
	}

	listModel := list.New(callers, list.NewDefaultDelegate(), 25, 50)
	listModel.Title = "Top-Level Callers"
	listModel.SetFilteringEnabled(true)
	vp := viewport.New(60, 40)
	vp.SetContent(renderPaths(pathsBycaller[callers[0].FilterValue()], target))

	return Model{
		callers:    listModel,
		paths:      vp,
		callerData: pathsBycaller,
		target:     target,
	}
}
