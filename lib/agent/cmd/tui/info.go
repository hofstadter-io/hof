package tui

import (
	"fmt"
	"maps"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type infoModel struct {
	root   *Model
	keymap infoKeymap

	// model specific
	Tabs      []string
	activeTab int
	viewport  viewport.Model
}

func initialInfoModel(root *Model) *infoModel {
	width, height := root.subwidth, root.subheight
	tabs := []string{"info", "msgs", "files", "state"}
	vh := height - 6
	vp := viewport.New(width, vh)

	return &infoModel{
		root:     root,
		keymap:   infoKeymapDefaults,
		Tabs:     tabs,
		viewport: vp,
	}
}

func (m *infoModel) Init() tea.Cmd {
	return nil
}

func (m *infoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.viewport.Height = m.root.subheight - 6
	m.viewport.Width = m.root.subwidth - 2

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
	if m.root.session == nil {
		content = "nil session"
	} else {
		name := m.Tabs[m.activeTab]
		switch name {
		case "info":
			content = m.renderTabInfo()
		case "files":
			content = m.renderTabFiles()
		case "state":
			content = m.renderTabState()
		case "msgs":
			content = m.renderTabMessages()
		}
	}
	m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(content))
	// ws := windowStyle.Width((width - windowStyle.GetHorizontalFrameSize()))
	// return ws.Render(content)
}

const infoFmt = `
id:     %s
title:  %s
events: %d
state:  %d

agent:  %s
model:  %s
envrn:  %s
`

func (m *infoModel) renderTabInfo() string {
	id := m.root.session.ID()
	state := maps.Collect(m.root.session.State().All())
	numEvents := m.root.session.Events().Len()
	numState := len(state)
	title := state["title"]
	agent := state["agent"]
	model := state["model"]
	envName := state["envName"]
	var b strings.Builder
	fmt.Fprintf(&b, strings.TrimSpace(infoFmt), id, title, numEvents, numState, agent, model, envName)
	return b.String()
}

func (m *infoModel) renderTabFiles() string {
	return "session files, tbd..."
}

func (m *infoModel) renderTabState() string {
	state := maps.Collect(m.root.session.State().All())
	pairs := make([]KVPair, 0, len(state))
	for k, v := range state {
		pairs = append(pairs, KVPair{Key: k, Value: v})
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Key < pairs[j].Key
	})

	var b strings.Builder
	for i, p := range pairs {
		fmt.Fprintf(&b, "%-3d%v\n", i, p.Key)
	}
	return b.String()
}

func (m *infoModel) renderTabMessages() string {

	msgs := renderMessages(m.root.subwidth, m.root.session)
	return strings.Join(msgs, "\n") + "\n\n"
}

type KVPair struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}
