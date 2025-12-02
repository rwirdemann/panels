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
	root   *panels.Panel
	panels map[int]*panels.Panel
	focus  int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.root.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl-c", "q":
			return m, tea.Quit
			// case "tab":
			// 	m.panels[m.focus].Blur()
			// 	if m.focus == 7 {
			// 		m.focus = 0
			// 	} else {
			// 		m.focus += 1
			// 	}
			// 	m.panels[m.focus].Focus()
			// 	return m, nil
		}

	}
	return m, nil
}

func (m model) View() string {
	return m.root.View(m, m.width, m.height)
}

func render(m tea.Model, panelID int, w, h int) string {
	return "Press 'v' for vertical split"
}

func main() {
	rootPanel := panels.NewPanel(10, panels.LayoutDirectionVertical, 100)
	m := model{root: rootPanel, panels: make(map[int]*panels.Panel)}

	row1 := panels.NewPanel(20, panels.LayoutDirectionHorizontal, 50)
	rootPanel.Append(row1)
	for i := range 4 {
		p := panels.NewPanel(i, panels.LayoutDirectionHorizontal, 25).WithBorder()
		row1.Append(p)
		m.panels[i] = p
	}

	row2 := panels.NewPanel(30, panels.LayoutDirectionHorizontal, 50)
	rootPanel.Append(row2)
	for i := 4; i < 8; i++ {
		p := panels.NewPanel(i, panels.LayoutDirectionHorizontal, 25).WithBorder()
		row2.Append(p)
		m.panels[i] = p
	}

	m.focus = 0
	m.panels[0].Focus()

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("there's been an error: %v", err)
		os.Exit(1)
	}
}
