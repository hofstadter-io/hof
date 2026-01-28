package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type infoModel struct {
	root   *Model
	keymap infoKeymap
	help   help.Model

	// sizing
	width  int
	height int

	// model specific
	Tabs      []string
	activeTab int
	viewport  viewport.Model

	// other common fields
	err error
}

func initialInfoModel(root *Model, width, height int) *infoModel {
	tabs := []string{"info", "state", "msgs", "Eye Shadow", "Mascara", "Foundation"}
	vh := height - 6
	vp := viewport.New(width-2, vh)

	return &infoModel{
		root:     root,
		keymap:   infoKeymapDefaults,
		help:     help.New(),
		width:    width,
		height:   height,
		Tabs:     tabs,
		viewport: vp,
	}
}

func (m *infoModel) Init() tea.Cmd {
	return nil
}

func (m *infoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.viewport.Height = m.height - 6
	m.viewport.Width = m.width - 2

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.renderTab()
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.back):
			m.root.clearSession()
			m.root.updateCurrName("list")
			cmd := m.root.runTick()
			return m, cmd
		case key.Matches(msg, m.keymap.chat):
			m.root.updateCurrName("chat")
			m.root.updateRootTitle()
			m.root.chat.updateMessagesFromEvents()
			cmd := m.root.runTick()
			return m, cmd
		case key.Matches(msg, m.keymap.next):
			m.activeTab = min(m.activeTab+1, len(m.Tabs)-1)
			m.renderTab()
			return m, nil
		case key.Matches(msg, m.keymap.prev):
			m.activeTab = max(m.activeTab-1, 0)
			m.renderTab()
			return m, nil

		default:
			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

var (
	docStyle         = lipgloss.NewStyle()
	highlightText    = lipgloss.AdaptiveColor{Light: "#a77ef8", Dark: "#a083fa"}
	highlightColor   = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(0, 1)
	activeTabStyle   = inactiveTabStyle.Foreground(highlightText).BorderForeground(highlightColor)
	windowStyle      = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(1, 1).Align(lipgloss.Left).Border(lipgloss.NormalBorder())
)

func (m *infoModel) View() string {
	doc := strings.Builder{}

	var renderedTabs []string

	for i, t := range m.Tabs {
		var style lipgloss.Style
		if i == m.activeTab {
			style = activeTabStyle
		} else {
			style = inactiveTabStyle
		}
		border, _, _, _, _ := style.GetBorder()
		// if isFirst && isActive {
		// 	border.BottomLeft = "│"
		// } else if isFirst && !isActive {
		// 	border.BottomLeft = "├"
		// } else if isLast && isActive {
		// 	border.BottomRight = "│"
		// } else if isLast && !isActive {
		// 	border.BottomRight = "┴"
		// }
		style = style.Border(border).UnsetBorderBottom()
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	// fstr := fmt.Sprintf("%%%ds\n", m.width)
	doc.WriteString(row)
	doc.WriteString("\n")
	doc.WriteString(windowStyle.Render(m.viewport.View()))
	return docStyle.Render(doc.String())
}

func (m *infoModel) renderTab() {
	content := "ipsum lorum"
	name := m.Tabs[m.activeTab]
	switch name {
	case "info":
		content = "session info..."
	case "state":
		content = "session state..."
	case "msgs":
		content = m.renderTabMessages()
	}
	m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(content))
	// ws := windowStyle.Width((width - windowStyle.GetHorizontalFrameSize()))
	// return ws.Render(content)
}

func (m *infoModel) renderTabMessages() string {

	msgs := renderMessages(m.width, m.root.session)
	return strings.Join(msgs, "\n") + "\n\n"
}
