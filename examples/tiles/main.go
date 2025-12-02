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
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		fmt.Printf("WindowSize: %v", msg)
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

func render(m tea.Model, panelID int, w, h int) string {
	return fmt.Sprintf("Tile %d", panelID)
}

func main() {
	rootPanel := panels.NewPanel(1, panels.LayoutDirectionVertical, 100).
		WithBorder().
		WithContent(render)
	m := model{panel: rootPanel}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("there's been an error: %v", err)
		os.Exit(1)
	}
}
