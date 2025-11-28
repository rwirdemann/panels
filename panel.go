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
	renderContent   func(m tea.Model, w, h int) string
}

func NewPanel(layout LayoutDirection, hasBorder bool, hasHelp bool, ratio float32, renderContent func(m tea.Model, w, h int) string) *Panel {
	return &Panel{layoutDirection: layout, hasBorder: hasBorder, hasHelp: hasHelp, ratio: ratio, renderContent: renderContent}
}

func (p *Panel) Append(panel *Panel) {
	p.children = append(p.children, panel)
}

func (p *Panel) distributeHorizontally(width int) {
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

func (p *Panel) distributeVertically(height int) {
	for i, child := range p.children {
		if child.hasBorder {
			p.children[i].height = height - 2
		} else {
			p.children[i].height = height
		}
	}
}

func (p *Panel) View(m tea.Model, parentWith, parentHeight int) string {
	if len(p.children) > 0 {

		if p.layoutDirection == LayoutDirectionHorizontal {
			p.distributeHorizontally(parentWith)
			p.distributeVertically(parentHeight)
			var children []string
			for _, c := range p.children {
				children = append(children, c.View(m, c.width, c.height))
			}
			return lipgloss.JoinHorizontal(lipgloss.Top, children...)
		}

		if p.layoutDirection == LayoutDirectionVertical {
			p.width = parentWith
			p.height = parentHeight

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
		content = content + p.renderContent(m, p.width-h, p.height-v)
	}
	return style.Render(content)
}
