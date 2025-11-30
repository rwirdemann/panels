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
	Name            string
	width           int
	height          int
	hasBorder       bool
	children        []*Panel
	layoutDirection LayoutDirection
	ratio           float32
	hasHelp         bool
	renderContent   func(m tea.Model, name string, w, h int) string
}

func NewPanel(layout LayoutDirection, hasBorder bool, hasHelp bool, ratio float32) *Panel {
	return &Panel{layoutDirection: layout, hasBorder: hasBorder, hasHelp: hasHelp, ratio: ratio}
}

func (p *Panel) WithContent(f func(m tea.Model, name string, w, h int) string) *Panel {
	p.renderContent = f
	return p
}

func (p *Panel) Append(panel *Panel) {
	p.children = append(p.children, panel)
}

func (p *Panel) View(m tea.Model, parentWith, parentHeight int) string {
	if len(p.children) > 0 {

		if p.layoutDirection == LayoutDirectionHorizontal {
			totalUsed := 0
			for i, child := range p.children {
				childWidth := int(float32(parentWith) * child.ratio)
				if child.hasBorder {
					p.children[i].width = childWidth - 2
					p.children[i].height = parentHeight - 2
				} else {
					p.children[i].width = childWidth
					p.children[i].height = parentHeight
				}
				totalUsed += childWidth
			}
			if totalUsed < parentWith && len(p.children) > 0 {
				p.children[len(p.children)-1].width += parentWith - totalUsed
			}
			var children []string
			for _, c := range p.children {
				children = append(children, c.View(m, c.width, c.height))
			}
			return lipgloss.JoinHorizontal(lipgloss.Top, children...)
		}

		if p.layoutDirection == LayoutDirectionVertical {
			totalUsed := 0
			for i, c := range p.children {
				height := int(float32(parentHeight) * c.ratio)
				totalUsed += height
				if c.hasBorder {
					p.children[i].width = parentWith - 2
					p.children[i].height = height - 2
				} else {
					p.children[i].width = parentWith
					p.children[i].height = height
				}
			}

			if totalUsed < parentHeight && len(p.children) > 0 {
				p.children[len(p.children)-1].height += parentHeight - totalUsed
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
		style = style.Border(lipgloss.NormalBorder())
	}
	content := ""
	if p.renderContent != nil {
		h, v := style.GetFrameSize()
		content = content + p.renderContent(m, p.Name, p.width-h, p.height-v)
	}
	return style.Render(content)
}
