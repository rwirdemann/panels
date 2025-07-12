package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	width  int
	height int
)

const (
	LayoutDirectionHorizontal = iota
	LayoutDirectionVertial
)

type Panel struct {
	name            string
	width           int
	height          int
	hasBorder       bool
	children        []Panel
	layoutDirection int
}

func (p Panel) Resize(width int, height int) Panel {
	p.height = height - 2
	p.width = width
	return p
}

func distribute(width int, panels []Panel) []int {
	if len(panels) == 0 {
		return []int{}
	}

	panelWidth := width / len(panels)
	widths := make([]int, len(panels))
	totalUsed := 0
	for i, panel := range panels {
		if panel.hasBorder {
			widths[i] = panelWidth - 2
			totalUsed += widths[i] + 2
		} else {
			widths[i] = panelWidth
			totalUsed += widths[i]
		}
	}

	if totalUsed < width && len(panels) > 0 {
		widths[len(panels)-1] += width - totalUsed
	}

	return widths
}

func (p Panel) View() string {
	if len(p.children) > 0 {
		var children []string
		widths := distribute(width, p.children)
		for i, c := range p.children {
			c = c.Resize(widths[i], height)
			children = append(children, c.View())
		}
		return lipgloss.JoinHorizontal(lipgloss.Top, children...)
	}

	style := lipgloss.NewStyle().Height(p.height).Width(p.width)
	if p.hasBorder {
		style = style.Border(lipgloss.NormalBorder())
	}
	return style.Render(fmt.Sprintf("Window: %dx%d\n Panel: %dx%d", width, height, p.width, p.height))
}

type model struct {
	panel Panel
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		width = msg.Width
		height = msg.Height
		m.panel = m.panel.Resize(width, height)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl-c", "q":
			return m, tea.Quit
		}

	}
	return m, nil
}

func (m model) View() string {
	return m.panel.View()
}

func main() {
	m := model{panel: Panel{hasBorder: true}}
	m.panel.children = append(m.panel.children, Panel{name: "Panel 1", hasBorder: true},
		Panel{name: "Panel 1", hasBorder: true})
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
