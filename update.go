package main

import tea "github.com/charmbracelet/bubbletea"

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Update list selection
		m.callers, cmd = m.callers.Update(msg)
		selected := m.callers.SelectedItem().(listItem)
		m.paths.SetContent(renderPaths(m.callerData[string(selected)], m.target))
	}

	return m, cmd
}
