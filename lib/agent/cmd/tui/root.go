package tui

import (
	"fmt"
	"maps"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"google.golang.org/adk/session"

	aruntime "github.com/hofstadter-io/hof/lib/agent/runtime"
	"github.com/hofstadter-io/hof/lib/agent/runtime/handlers/common"
	"github.com/hofstadter-io/hof/lib/runtime"
)

const gap = "\n\n"

type (
	errMsg error
)

type Model struct {
	// the core runtimes we can work with
	R  *runtime.Runtime
	AR *aruntime.Runtime

	// sizing
	width  int
	height int

	// keymap & help
	keymap rootKeymap
	help   help.Model

	// curr
	curr      tea.Model
	currIdx   int
	currName  string
	currStyle lipgloss.Style

	session   session.Session
	currSid   string
	currTitle string

	// views
	// main mainModel
	list *listModel
	info *infoModel
	chat *chatModel

	// other fields
	msg string
	err error
}

var views = []string{"list", "chat", "info"}

func InitialModel(R *runtime.Runtime, AR *aruntime.Runtime) *Model {
	w := 50
	h := 25

	cn := "list"
	cs := lipgloss.NewStyle().Foreground(lipgloss.Color("5"))

	m := &Model{
		R:      R,
		AR:     AR,
		width:  w,
		height: h,
		help:   help.New(),
		keymap: rootKeymapDefaults,
		err:    nil,

		// curr:      sess,
		currName:  cn,
		currStyle: cs,
		currIdx:   0,
		// sess:      sess,
	}

	sess := initialSessionsModel(m, w, h)
	m.list = sess
	m.curr = sess

	return m
}

type TickMsg time.Time

// Send a message every second.
func (m *Model) runTick() tea.Cmd {
	return tea.Every(time.Millisecond*300, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m *Model) Init() tea.Cmd {
	return m.runTick()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		mainCmd tea.Cmd
		tickCmd tea.Cmd
	)

	if msg == nil {
		return m, tea.Sequence(mainCmd)
	}

	// this is getting message, especially the key handling, which has been resistant to nested handling and passing around control
	switch msg := msg.(type) {
	// This is needed so that ticks act like intervals and we have a polling sync like mechanism
	case TickMsg:
		// Return your Every command again to loop.
		// tickCmd = m.runTick()

	case tea.WindowSizeMsg:
		m.height = msg.Height - 1
		m.width = msg.Width - 1
		m.currDims()
		newCurr, cmd := m.curr.Update(msg)
		m.curr = newCurr
		return m, tea.Sequence(cmd, mainCmd)

	case tea.KeyMsg:
		// handle normally
		switch {
		case key.Matches(msg, m.keymap.help):
			if m.currName != "chat" || !m.chat.textarea.Focused() {
				show := !m.help.ShowAll
				m.help.ShowAll = show
				if m.list != nil {
					m.list.help.ShowAll = show
				}
				if m.info != nil {
					m.info.help.ShowAll = show
				}
				if m.chat != nil {
					m.chat.help.ShowAll = show
				}
			}
		case key.Matches(msg, m.keymap.quit):
			// special handling of for chat
			var match bool
			// find situations we want to pass 'q' through and not quit
			if (m.currName == "chat" && m.chat != nil && m.chat.textarea.Focused()) || (m.currName == "list" && m.list.input.Focused()) {
				switch keypress := msg.String(); keypress {
				case "q":
					match = true
				}
			}
			if !match {
				return m, tea.Quit
			}
		}

		m.currDims()
		_, cmd := m.curr.Update(msg)
		return m, tea.Sequence(cmd, mainCmd, tickCmd)

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		switch m.currName {
		case "list":
			m.list.err = msg
		case "info":
			m.info.err = msg
		case "chat":
			m.chat.err = msg
		}
		return m, nil
	}

	return m, tea.Sequence(mainCmd)
}

func (m *Model) View() string {
	doc := strings.Builder{}

	// header row
	var first string
	switch m.currName {
	case "list":
		first = "sessions"
	case "chat", "info":
		if m.currTitle != "" {
			first = m.currTitle
		} else if m.currSid != "" {
			first = m.currSid
		}

	default:
		first = "error, unknown view: " + m.currName
	}

	first = m.currStyle.Render(first)

	help := m.help.View(m.keymap)
	var subhelp string
	switch m.currName {
	case "list":
		subhelp = m.list.help.View(m.list.keymap)
	case "info":
		subhelp = m.info.help.View(m.info.keymap)
	case "chat":
		subhelp = m.chat.help.View(m.chat.keymap)
	}
	header := lipgloss.JoinHorizontal(lipgloss.Top, subhelp, " | ", help, " | ", first)
	doc.WriteString(header)
	doc.WriteString("\n")

	// main view
	curr := m.curr.View()
	doc.WriteString(curr)

	return doc.String()
}

func (m *Model) updateCurr(i int) {
	// update index
	lv := len(views)
	i = (i + lv) % lv
	m.currIdx = i

	// update name
	cn := views[i]
	m.currName = cn

	// update curr
	switch cn {
	case "info":
		if m.info == nil {
			m.info = initialInfoModel(m, m.width, m.height)
		}
		m.curr = m.info
	case "list":
		if m.list == nil {
			m.list = initialSessionsModel(m, m.width, m.height)
		}
		m.curr = m.list
	case "chat":
		if m.chat == nil {
			m.chat = initialChatModel(m, m.width, m.height)
		}
		m.curr = m.chat
		m.chat.textarea.Focus()
	}
}

func (m *Model) updateCurrName(view string) {
	// find name
	i := -1
	for j, v := range views {
		if v == view {
			i = j
			break
		}
	}
	if i == -1 {
		m.msg = "unknown view: " + view
	}

	m.updateCurr(i)
}

func (m *Model) clearSession() {
	m.err = nil
	m.msg = ""
	m.currSid = ""
	m.currTitle = ""
}

func (m *Model) currDims() {
	var hcnt, ccnt, icnt, scnt int
	htxt := m.help.View(m.keymap)
	hcnt = strings.Count(htxt, "\n") + 1

	if m.chat != nil && m.currName == "chat" {
		htxt = m.chat.help.View(m.chat.keymap)
		ccnt = strings.Count(htxt, "\n") + 1
	}
	if m.info != nil && m.currName == "info" {
		htxt = m.info.help.View(m.info.keymap)
		icnt = strings.Count(htxt, "\n") + 1
	}
	if m.list != nil && m.currName == "list" {
		htxt = m.list.help.View(m.list.keymap)
		scnt = strings.Count(htxt, "\n") + 1
	}
	hcnt = max(hcnt, ccnt, icnt, scnt)

	vh := m.height - hcnt + 1

	if m.list != nil {
		m.list.width = m.width
		m.list.height = vh
	}
	if m.info != nil {
		m.info.width = m.width
		m.info.height = vh
	}
	if m.chat != nil {
		m.chat.width = m.width
		m.chat.height = vh
	}
}

func (m *Model) loadSession(sid string) error {
	m.currSid = sid
	session, err := common.SessionGet(m.R, m.AR, m.currSid)
	if err != nil {
		m.msg = err.Error()
		return err
	} else {
		m.session = session
	}

	return nil
}

func (m *Model) updateRootTitle() {
	state := maps.Collect(m.session.State().All())
	numEvents := m.session.Events().Len()
	numState := len(state)
	title := state["title"]
	m.currTitle = fmt.Sprintf("%s  events:%d  state:%d", title, numEvents, numState)
}
