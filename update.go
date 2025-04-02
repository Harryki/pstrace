package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "[":
			m.paths.LineUp(m.paths.Height)
		case "]":
			m.paths.LineDown(m.paths.Height)
		case "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		_, v := docStyle.GetFrameSize()
		m.callers.SetSize(msg.Width, msg.Height-v)
		// viewport.Sync(m.paths)

		// m.paths.Width = msg.Width
		m.paths.Height = msg.Height - v - 1
	}

	var cmd tea.Cmd
	m.callers, cmd = m.callers.Update(msg) // ← always handle all msg types

	if selectedItem := m.callers.SelectedItem(); selectedItem != nil {
		selected := selectedItem.(listItem)
		m.paths.SetContent(renderPaths(m.callerData[string(selected)], m.target))
	}

	return m, cmd
}
