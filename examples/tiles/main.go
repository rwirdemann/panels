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
	return "Press 'v' for vertical split"
}

func main() {
	rootPanel := panels.NewPanel(10, panels.LayoutDirectionVertical, 100)
	m := model{panel: rootPanel}

	row1 := panels.NewPanel(20, panels.LayoutDirectionHorizontal, 50)
	rootPanel.Append(row1)
	for i := 21; i < 25; i++ {
		row1.Append(panels.NewPanel(i, panels.LayoutDirectionHorizontal, 25).WithBorder())
	}

	row2 := panels.NewPanel(30, panels.LayoutDirectionHorizontal, 50)
	rootPanel.Append(row2)
	for i := 31; i < 35; i++ {
		row2.Append(panels.NewPanel(i, panels.LayoutDirectionHorizontal, 25).WithBorder())
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("there's been an error: %v", err)
		os.Exit(1)
	}
}
