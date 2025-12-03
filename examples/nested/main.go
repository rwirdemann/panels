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

func renderPanel(m tea.Model, panelID int, w, h int) string {
	return fmt.Sprintf("Panel %d", panelID)
}

func main() {
	rootPanel := panels.NewPanel(1, 100).WithLayout(panels.LayoutDirectionHorizontal)
	m := model{panel: rootPanel}
	leftPanel := panels.NewPanel(2, 50).
		WithContent(renderPanel).
		WithBorder()
	leftPanel.Focus()
	rootPanel.Append(leftPanel)

	rightPanel := panels.NewPanel(3, 50).WithLayout(panels.LayoutDirectionVertical)
	rootPanel.Append(rightPanel)

	topPanel := panels.NewPanel(4, 50).
		WithContent(renderPanel).
		WithBorder()

	rightPanel.Append(topPanel)
	bottomPanel := panels.NewPanel(5, 50).WithLayout(panels.LayoutDirectionHorizontal)
	rightPanel.Append(bottomPanel)

	leftBottomPanel := panels.NewPanel(6, 50).
		WithContent(renderPanel).
		WithBorder()
	bottomPanel.Append(leftBottomPanel)

	rightBottomPanel := panels.NewPanel(7, 50).
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
