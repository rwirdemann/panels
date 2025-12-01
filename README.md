# Bubble Tea Panels

Bubble Tea Panels provides a framework for placing, layouting and resizing panels across the root window of a [Bubble Tea](https://github.com/charmbracelet/bubbletea) application.

## Basic Usage

### The Model

Add a root panel together with width and height to your Bubble Tea model.

```go
type model struct {
	root    *panels.Panel
	width   int
	height  int
}
```
  
### The Update Method  

```
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
  ...
  }
  ...
}

func (m model) View() string {
	return m.panel.View(m, m.width, m.height)
}

func renderPanel(m tea.Model, name string, w, h int) string {
	return name
}


func main() {
  rootPanel := panels.NewPanel(panels.LayoutDirectionHorizontal, true, false, 1.0)
  m := model{panel: rootPanel, list: list.New(items, list.NewDefaultDelegate(), 0, 0)}

  leftPanel := panels.NewPanel(panels.LayoutDirectionNone, true, false, 0.35).WithContent(renderPanel)
  rootPanel.Append(leftPanel)
  rightPanel := panels.NewPanel(panels.LayoutDirectionNone, true, false, 0.65).WithContent(renderPanel)
  rootPanel.Append(rightPanel)

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
```

## Horizontal Layout

```go
rootPanel := panels.NewPanel(panels.LayoutDirectionHorizontal, true, false, 1.0)
m := model{panel: rootPanel, list: list.New(items, list.NewDefaultDelegate(), 0, 0)}

leftPanel := panels.NewPanel(panels.LayoutDirectionNone, true, false, 0.35).WithContent(left)
rootPanel.Append(leftPanel)
rightPanel := panels.NewPanel(panels.LayoutDirectionNone, true, false, 0.65).WithContent(right)
rootPanel.Append(rightPanel)
┌───────────────────────────────────────┐
│ ┌────────────────┐ ┌────────────────┐ │
│ │ left           │ │ right          │ │
│ │                │ │                │ │
│ │                │ │                │ │
│ │                │ │                │ │
│ │                │ │                │ │
│ │                │ │                │ │
│ │                │ │                │ │
│ │                │ │                │ │
│ │                │ │                │ │
│ └────────────────┘ └────────────────┘ │
└───────────────────────────────────────┘
```
