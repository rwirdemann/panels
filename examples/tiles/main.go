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

func main() {
	rootPanel := panels.NewPanel(10, 100).WithLayout(panels.LayoutDirectionVertical)
	m := model{root: rootPanel, panels: make(map[int]*panels.Panel)}

	row1 := panels.NewPanel(20, 50).WithLayout(panels.LayoutDirectionHorizontal)
	rootPanel.Append(row1)
	for i := range 4 {
		p := panels.NewPanel(i, 25).WithBorder()
		row1.Append(p)
		m.panels[i] = p
	}

	row2 := panels.NewPanel(30, 50).WithLayout(panels.LayoutDirectionHorizontal)
	rootPanel.Append(row2)
	for i := 4; i < 8; i++ {
		p := panels.NewPanel(i, 25).WithBorder()
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
