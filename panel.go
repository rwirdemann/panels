package panels

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

type LayoutDirection int

const (
	LayoutDirectionNone       LayoutDirection = 0
	LayoutDirectionHorizontal LayoutDirection = 1
	LayoutDirectionVertical   LayoutDirection = 2
)

var PanelBorderColorFocus lipgloss.Color
var PanelBorderColor lipgloss.Color

func init() {
	isDark := termenv.HasDarkBackground()
	if isDark {
		PanelBorderColorFocus = "12"
		PanelBorderColor = "255"
	} else {
		PanelBorderColorFocus = "12"
		PanelBorderColor = "0"
	}
}

type Panel struct {
	ID              int
	Title           string
	width           int
	height          int
	hasBorder       bool
	children        []*Panel
	layoutDirection LayoutDirection
	ratio           int
	hasHelp         bool
	hasFocus        bool
	renderContent   func(m tea.Model, panelID int, w, h int) string
}

func NewPanel() *Panel {
	return &Panel{layoutDirection: LayoutDirectionNone}
}

func (p *Panel) WithId(id int) *Panel {
	p.ID = id
	return p
}

func (p *Panel) WithRatio(ratio int) *Panel {
	p.ratio = ratio
	return p
}

func (p *Panel) WithTitle(title string) *Panel {
	p.Title = title
	return p
}

func (p *Panel) WithContent(f func(m tea.Model, panelID int, w, h int) string) *Panel {
	p.renderContent = f
	return p
}

func (p *Panel) WithLayout(layout LayoutDirection) *Panel {
	p.layoutDirection = layout
	return p
}

func (p *Panel) WithBorder() *Panel {
	p.hasBorder = true
	return p
}

func (p *Panel) WithHelp() *Panel {
	p.hasHelp = true
	return p
}

func (p *Panel) Focus() {
	p.hasFocus = true
}

func (p *Panel) Blur() {
	p.hasFocus = false
}

func (p *Panel) Append(panel *Panel) {
	p.children = append(p.children, panel)
}

func (p *Panel) Update(msg tea.Msg) (*Panel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			// Find the currently focused panel
			var currentFocused *Panel
			p.walk(func(panel *Panel) bool {
				if panel.hasFocus {
					currentFocused = panel
					return false // Stop traversal
				}
				return true
			})

			if currentFocused != nil {
				currentFocused.Blur()
			}

			// Find the next panel to focus
			nextFocus := p.findNextFocusablePanel(currentFocused)

			if nextFocus != nil {
				nextFocus.Focus()
			} else if currentFocused == nil {
				// If no panel was focused initially, focus the first one
				p.walk(func(panel *Panel) bool {
					nextFocus = panel
					return false
				})
				if nextFocus != nil {
					nextFocus.Focus()
				}
			}
			return p, nil
		}
	}

	return p, nil
}

// findNextFocusablePanel determines the next panel to focus based on the
// current one. It implements a depth-first traversal logic and only considers
// leaf panels (panels without children) as focusable.
func (p *Panel) findNextFocusablePanel(currentFocused *Panel) *Panel {
	var focusablePanels []*Panel
	p.walk(func(panel *Panel) bool {
		if len(panel.children) == 0 {
			focusablePanels = append(focusablePanels, panel)
		}
		return true
	})

	if len(focusablePanels) == 0 {
		return nil
	}

	if currentFocused == nil {
		return focusablePanels[0]
	}

	for i, panel := range focusablePanels {
		if panel == currentFocused {
			if i < len(focusablePanels)-1 {
				return focusablePanels[i+1] // Return the next leaf panel
			} else {
				return focusablePanels[0] // Wrap around to the first leaf panel
			}
		}
	}

	// If the currently focused panel is not in the list of focusable panels
	// (e.g., it's a container that somehow got focus), return the first
	// focusable panel.
	return focusablePanels[0]
}

// walk performs a depth-first traversal of the panel tree, calling the visitor
// function for each panel. If the visitor returns false, the traversal stops.
func (p *Panel) walk(visitor func(*Panel) bool) bool {
	if !visitor(p) {
		return false
	}
	for _, child := range p.children {
		if !child.walk(visitor) {
			return false
		}
	}
	return true
}

func (p *Panel) View(m tea.Model, parentWidth, parentHeight int) string {
	// ignore first view call when windows size is still 0
	if parentWidth == 0 && parentHeight == 0 {
		return ""
	}

	// Set panel dimensions from parent if not already set by parent panel
	if p.width == 0 || p.height == 0 {
		if p.hasBorder {
			p.width = parentWidth - 2
			p.height = parentHeight - 2
		} else {
			p.width = parentWidth
			p.height = parentHeight
		}
	}

	if len(p.children) > 0 {

		if p.layoutDirection == LayoutDirectionHorizontal {
			totalWeight := 0
			for _, child := range p.children {
				totalWeight += child.ratio
			}

			if totalWeight > 0 {
				remainder := parentWidth
				for i, child := range p.children {
					childWidth := (parentWidth * child.ratio) / totalWeight
					if child.hasBorder {
						p.children[i].width = childWidth - 2
						p.children[i].height = parentHeight - 2
					} else {
						p.children[i].width = childWidth
						p.children[i].height = parentHeight
					}
					remainder -= childWidth
				}
				// Distribute remainder
				for i := 0; i < remainder; i++ {
					p.children[i%len(p.children)].width++
				}
			}

			var children []string
			for _, c := range p.children {
				children = append(children, c.View(m, c.width, c.height))
			}
			return lipgloss.JoinHorizontal(lipgloss.Top, children...)
		}

		if p.layoutDirection == LayoutDirectionVertical {
			totalWeight := 0
			for _, child := range p.children {
				totalWeight += child.ratio
			}

			if totalWeight > 0 {
				remainder := parentHeight
				for i, c := range p.children {
					height := (parentHeight * c.ratio) / totalWeight
					remainder -= height
					if c.hasBorder {
						p.children[i].width = parentWidth - 2
						p.children[i].height = height - 2
					} else {
						p.children[i].width = parentWidth
						p.children[i].height = height
					}
				}

				// Distribute remainder
				for i := 0; i < remainder; i++ {
					p.children[i%len(p.children)].height++
				}
			}

			var children []string
			for _, c := range p.children {
				children = append(children, c.View(m, c.width, c.height))
			}
			return lipgloss.JoinVertical(lipgloss.Top, children...)
		}
	}

	style := lipgloss.NewStyle().Height(p.height).Width(p.width)
	if p.hasBorder {
		style = style.Border(lipgloss.RoundedBorder())
	}
	if p.hasFocus {
		style = style.BorderForeground(PanelBorderColorFocus)
	}
	content := ""
	if p.renderContent != nil {
		h, v := style.GetFrameSize()
		content = content + p.renderContent(m, p.ID, p.width-h, p.height-v)
	}
	return style.Render(content)
}
