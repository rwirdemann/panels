package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rwirdemann/panels"
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	width   int
	height  int
	panel   *panels.Panel
	list    list.Model
	content string
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
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	// Update content with description of currently selected item
	if i, ok := m.list.SelectedItem().(item); ok {
		m.content = i.desc
	}

	return m, cmd
}

func (m model) View() string {
	return m.panel.View(m, m.width, m.height)
}

func listView(m tea.Model, panelID int, w, h int) string {
	model := m.(model)
	model.list.SetSize(w, h)
	return model.list.View()
}

func right(m tea.Model, panelID int, w, h int) string {
	model := m.(model)
	return model.content
}

func main() {
	items := []list.Item{
		item{title: "Raspberry Pi’s", desc: "I have ’em all over my house"},
		item{title: "Nutella", desc: "It's good on toast"},
		item{title: "Bitter melon", desc: "It cools you down"},
	}

	rootPanel := panels.NewPanel(1, panels.LayoutDirectionHorizontal, true, false, 1.0)
	m := model{panel: rootPanel, list: list.New(items, list.NewDefaultDelegate(), 0, 0)}

	leftPanel := panels.NewPanel(2, panels.LayoutDirectionNone, true, false, 0.35).
		WithContent(listView)
	rootPanel.Append(leftPanel)
	rightPanel := panels.NewPanel(3, panels.LayoutDirectionNone, true, false, 0.65).WithContent(right)
	rootPanel.Append(rightPanel)

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
