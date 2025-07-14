package panels

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type LayoutDirection int

const (
	LayoutDirectionHorizontal LayoutDirection = 1
	LayoutDirectionVertial    LayoutDirection = 2
)

var (
	Width  int
	Height int
)

type Panel struct {
	width           int
	height          int
	hasBorder       bool
	children        []*Panel
	layoutDirection LayoutDirection
	ratio           float32
	hasHelp         bool
	renderContent   func(m tea.Model, w, h int) string
}

func NewPanel(hasBorder bool, hasHelp bool, ratio float32, renderContent func(m tea.Model, w, h int) string) *Panel {
	return &Panel{hasBorder: hasBorder, hasHelp: hasHelp, ratio: ratio, renderContent: renderContent}
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
	v := 2
	if p.hasHelp {
		v = 3
	}
	for i := range p.children {
		p.children[i].height = height - v
	}
}

func (p *Panel) View(m tea.Model) string {
	if len(p.children) > 0 {
		p.distributeHorizontally(Width)
		p.distributeVertically(Height)
		var children []string
		for _, c := range p.children {
			children = append(children, c.View(m))
		}
		return lipgloss.JoinHorizontal(lipgloss.Top, children...)
	}

	style := lipgloss.NewStyle().Height(p.height).Width(p.width)
	if p.hasBorder {
		style = style.Border(lipgloss.NormalBorder())
	}
	content := ""
	if p.renderContent != nil {
		h, v := style.GetFrameSize()
		content = p.renderContent(m, p.width-h, p.height-v)
	}
	return style.Render(content)
}
