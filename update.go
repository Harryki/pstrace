package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.callers, cmd = m.callers.Update(msg) // ← always handle all msg types

	if selectedItem := m.callers.SelectedItem(); selectedItem != nil {
		selected := selectedItem.(listItem)
		m.paths.SetContent(renderPaths(m.callerData[string(selected)], m.target))
	}

	return m, cmd
}
