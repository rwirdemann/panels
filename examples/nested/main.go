package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rwirdemann/panels"
)

type model struct {
	panel  *panels.Panel
	width  int
	height int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func renderPanel(m tea.Model, panelID int, w, h int) string {
	return fmt.Sprintf("Panel %d", panelID)
}

func main() {
	rootPanel := panels.NewPanel(1, panels.LayoutDirectionHorizontal, 1.0)
	m := model{panel: rootPanel}
	leftPanel := panels.NewPanel(2, panels.LayoutDirectionNone, 0.50).
		WithContent(renderPanel).
		WithBorder()
	rootPanel.Append(leftPanel)

	rightPanel := panels.NewPanel(3, panels.LayoutDirectionVertical, 0.50)
	rootPanel.Append(rightPanel)

	topPanel := panels.NewPanel(4, panels.LayoutDirectionNone, 0.50).
		WithContent(renderPanel).
		WithBorder()

	rightPanel.Append(topPanel)
	bottomPanel := panels.NewPanel(5, panels.LayoutDirectionHorizontal, 0.50)
	rightPanel.Append(bottomPanel)

	leftBottomPanel := panels.NewPanel(6, panels.LayoutDirectionNone, 0.50).
		WithContent(renderPanel).
		WithBorder()
	bottomPanel.Append(leftBottomPanel)

	rightBottomPanel := panels.NewPanel(7, panels.LayoutDirectionNone, .50).
		WithContent(renderPanel).
		WithBorder()
	rightBottomPanel.Title = "bottom right"
	bottomPanel.Append(rightBottomPanel)

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
