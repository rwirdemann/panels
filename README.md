# Bubble Tea Nestile

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

Make sure to add a case for `tea.WindowSizeMsg` and save the new width and height in your model.

```go
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
    ...
  }
  ...
}
```

### The View Method

Forward the method call to your root panel passing the model together with its width and height.

```go
func (m model) View() string {
	return m.root.View(m, m.width, m.height)
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
