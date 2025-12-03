# Bubble Tea Nestile

Bubble Tea Panels provides a framework for placing, layouting and resizing panels across the root window of a [Bubble Tea](https://github.com/charmbracelet/bubbletea) application.

```
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

Forward Bubble Tea's `View` call to your root panel's `View` method passing the model together with its width and height.

```go
func (m model) View() string {
	return m.root.View(m, m.width, m.height)
}
```

### The main function

Create your model by assigning a root panel that consumes the full size of your main window (ratio = 1.0). Every panel that is added to your root panel will be layouted according the root panels base layout and the ratio of the nested panel. 

```go
func main() {
	rootPanel := panels.NewPanel().WithId(1).
		WithRatio(100).
		WithLayout(panels.LayoutDirectionVertical)
	m := model{panel: rootPanel}

	topPanel := panels.NewPanel().WithId(2).
		WithRatio(50).
		WithContent(top).
		WithBorder()
	rootPanel.Append(topPanel)

	bottomPanel := panels.NewPanel().WithId(3).
		WithRatio(50).
		WithContent(bottom).
		WithBorder()
	rootPanel.Append(bottomPanel)

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("there's been an error: %v", err)
		os.Exit(1)
	}
}
```

### The missing part - Render functions

Each panel with content needs a render function that is responsible for rendering the panels content.

```go
func top(m tea.Model, panelID int, w, h int) string {
	return "top"
}
```

You can either add specific render functions to each panel using `WithContent` or stick to a generic function used by all panels. You can use the provided paneID to determine wich panel needs to be rendered.
