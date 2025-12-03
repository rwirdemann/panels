package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rwirdemann/panels"
)

type model struct {
	width  int
	height int
	panel  *panels.Panel
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.panel.Update(msg)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl-c", "q":
			return m, tea.Quit
		}

	}
	return m, nil
}

func (m model) View() string {
	return m.panel.View(m, m.width, m.height)
}

func top(m tea.Model, panelID int, w, h int) string {
	return "top"
}

func bottom(m tea.Model, panelID int, w, h int) string {
	return "bottom"
}

func main() {
	rootPanel := panels.NewPanel().WithId(1).WithRatio(100).WithLayout(panels.LayoutDirectionVertical)
	m := model{panel: rootPanel}

	topPanel := panels.NewPanel().WithId(2).WithRatio(50).
		WithContent(top).
		WithBorder()
	topPanel.Focus()
	rootPanel.Append(topPanel)

	bottomPanel := panels.NewPanel().WithId(3).WithRatio(50).
		WithContent(bottom).
		WithBorder()
	rootPanel.Append(bottomPanel)

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("there's been an error: %v", err)
		os.Exit(1)
	}
}
