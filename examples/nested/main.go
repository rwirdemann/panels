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

func renderPanel(m tea.Model, name string, w, h int) string {
	return name
}

func main() {
	rootPanel := panels.NewPanel(panels.LayoutDirectionHorizontal, true, false, 1.0)
	m := model{panel: rootPanel}
	leftPanel := panels.NewPanel(panels.LayoutDirectionNone, true, false, 0.50).WithContent(renderPanel)
	leftPanel.Name = "left"
	rootPanel.Append(leftPanel)

	rightPanel := panels.NewPanel(panels.LayoutDirectionVertical, false, false, 0.50).WithContent(renderPanel)
	rightPanel.Name = "right"
	rootPanel.Append(rightPanel)

	topPanel := panels.NewPanel(panels.LayoutDirectionNone, true, false, 0.50).WithContent(renderPanel)
	topPanel.Name = "top"

	rightPanel.Append(topPanel)
	bottomPanel := panels.NewPanel(panels.LayoutDirectionHorizontal, false, false, 0.50).WithContent(renderPanel)
	bottomPanel.Name = "bottom"
	rightPanel.Append(bottomPanel)

	leftBottomPanel := panels.NewPanel(panels.LayoutDirectionNone, true, false, 0.50).WithContent(renderPanel)
	leftBottomPanel.Name = "bottom left"
	bottomPanel.Append(leftBottomPanel)

	rightBottomPanel := panels.NewPanel(panels.LayoutDirectionNone, true, false, 0.50).WithContent(renderPanel)
	rightBottomPanel.Name = "bottom right"
	bottomPanel.Append(rightBottomPanel)

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
