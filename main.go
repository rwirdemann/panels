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
	width           int
	height          int
	hasBorder       bool
	children        []Panel
	layoutDirection int
	ratio           float32
}

func (p Panel) Resize(width int, height int) Panel {
	p.height = height - 2
	p.width = width
	return p
}

func (p Panel) distributeHorizontally(width int) {
	if len(p.children) == 0 {
		return
	}

	totalUsed := 0
	for i, panel := range p.children {
		panelWidth := int(float32(width) * panel.ratio)
		if panel.hasBorder {
			p.children[i].width = panelWidth - 2
		} else {
			p.children[i].width = panelWidth
		}
		totalUsed += panelWidth
	}

	if totalUsed < width && len(p.children) > 0 {
		p.children[len(p.children)-1].width += width - totalUsed
	}
}

func (p Panel) distributeVertically(height int) {
	for i := range p.children {
		p.children[i].height = height - 2
	}
}

func (p Panel) View() string {
	if len(p.children) > 0 {
		p.distributeHorizontally(width)
		p.distributeVertically(height)
		var children []string
		for _, c := range p.children {
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
	m.panel.children = append(m.panel.children,
		Panel{hasBorder: true, ratio: 0.35},
		Panel{hasBorder: true, ratio: 0.45},
		Panel{hasBorder: true, ratio: 0.20})
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
