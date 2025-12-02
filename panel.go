package panels

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type LayoutDirection int

const (
	LayoutDirectionNone       LayoutDirection = 0
	LayoutDirectionHorizontal LayoutDirection = 1
	LayoutDirectionVertical   LayoutDirection = 2
)

type Panel struct {
	ID              int
	Title           string
	width           int
	height          int
	hasBorder       bool
	children        []*Panel
	layoutDirection LayoutDirection
	weight          int
	hasHelp         bool
	renderContent   func(m tea.Model, panelID int, w, h int) string
}

func NewPanel(id int, layout LayoutDirection, weight int) *Panel {
	return &Panel{ID: id, layoutDirection: layout, weight: weight}
}

func (p *Panel) WithContent(f func(m tea.Model, panelID int, w, h int) string) *Panel {
	p.renderContent = f
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

func (p *Panel) Append(panel *Panel) {
	p.children = append(p.children, panel)
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
				totalWeight += child.weight
			}

			if totalWeight > 0 {
				remainder := parentWidth
				for i, child := range p.children {
					childWidth := (parentWidth * child.weight) / totalWeight
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
				totalWeight += child.weight
			}

			if totalWeight > 0 {
				remainder := parentHeight
				for i, c := range p.children {
					height := (parentHeight * c.weight) / totalWeight
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
	content := ""
	if p.renderContent != nil {
		h, v := style.GetFrameSize()
		content = content + p.renderContent(m, p.ID, p.width-h, p.height-v)
	}
	return style.Render(content)
}
