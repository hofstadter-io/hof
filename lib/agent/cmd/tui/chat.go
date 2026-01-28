package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type chatModel struct {
	root   *Model
	keymap chatKeymap
	help   help.Model

	// sizing
	width  int
	height int

	// model specific
	chatStyle  lipgloss.Style
	viewport   viewport.Model
	messages   []string
	textarea   textarea.Model
	userStyle  lipgloss.Style
	agentStyle lipgloss.Style
	funcStyle  lipgloss.Style
	markglam   *glamour.TermRenderer

	// other common fields
	err error
}

func initialChatModel(root *Model, width, height int) *chatModel {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 8192

	ta.SetWidth(width)
	ta.SetHeight(7)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false
	vh := height - ta.Height() - lipgloss.Height(gap)*2

	vp := viewport.New(width, vh)
	vp.SetContent(`Welcome to the chat room!
Type a message and press Enter to send.`)

	// ta.KeyMap.InsertNewline.SetEnabled(false)

	glam, _ := glamour.NewTermRenderer(
		glamour.WithStylePath("dark"),
		glamour.WithWordWrap(width),
	)

	return &chatModel{
		chatStyle:  lipgloss.NewStyle().BorderTopForeground(lipgloss.Color("4")).BorderTop(true),
		root:       root,
		keymap:     chatKeymapDefaults,
		help:       help.New(),
		width:      width,
		height:     height,
		textarea:   ta,
		messages:   []string{},
		viewport:   vp,
		userStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("4")),
		agentStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("2")),
		funcStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("3")),
		markglam:   glam,
		err:        nil,
	}
}

func (m *chatModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m *chatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.viewport.Height = m.height - m.textarea.Height() - lipgloss.Height(gap)*2
	m.viewport.Width = m.width
	m.textarea.SetWidth(m.width)

	if m.textarea.Focused() {
		m.textarea, tiCmd = m.textarea.Update(msg)
	} else {
		m.viewport, vpCmd = m.viewport.Update(msg)
	}

	glam, _ := glamour.NewTermRenderer(
		glamour.WithStylePath("dark"),
		glamour.WithWordWrap(m.viewport.Width),
	)
	m.markglam = glam
	// if len(m.messages) > 0 {
	// 	// Wrap content before setting it.
	// 	m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
	// }

	if m.root.session != nil {
		m.root.updateRootTitle()

		// if m.root.session.Events().Len() != len(m.messages) {
		m.updateMessagesFromEvents()
		// m.viewport.GotoBottom()
		// }
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// m.root.msg = msg.String()
		switch {
		case key.Matches(msg, m.keymap.back):
			if m.textarea.Focused() {
				m.textarea.Blur()
			} else {
				m.root.clearSession()
				m.root.updateCurrName("list")
			}
		case key.Matches(msg, m.keymap.send):
			if m.textarea.Focused() {
				// TODO, actually send the message
				content := m.textarea.Value()
				lines := strings.Split(content, "\n")
				m.root.msg = lines[0]
				m.textarea.Reset()
			} else {
				if !m.textarea.Focused() {
					m.textarea.Focus()
				}
			}
		case key.Matches(msg, m.keymap.focus):
			m.textarea.Focus()
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m *chatModel) View() string {
	s := fmt.Sprintf(
		"%s%s%s%s%s",
		m.viewport.View(),
		gap,
		fmt.Sprintln(m.root.msg),
		gap,
		m.textarea.View(),
	)

	ws := m.chatStyle.Width((m.width - windowStyle.GetHorizontalFrameSize()))
	return ws.Render(s)
}

func (m *chatModel) refresh() {
	m.root.loadSession(m.root.currSid)
	m.root.updateRootTitle()
	m.updateMessagesFromEvents()
	m.viewport.GotoBottom()
}

func (m *chatModel) updateMessagesFromEvents() {
	m.messages = renderMessages(m.width, m.root.session)
	content := strings.Join(m.messages, "\n") + "\n\n\n\n"
	m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(content))
}
