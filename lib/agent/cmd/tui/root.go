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
	"github.com/hofstadter-io/hof/lib/agent/services/environ"
	"github.com/hofstadter-io/hof/lib/consts"
	"github.com/hofstadter-io/hof/lib/runtime"
)

type (
	errMsg error
)

type ViewModel interface {
	init(root *Model)
}

type Model struct {
	// the core runtimes we can work with
	R  *runtime.Runtime
	AR *aruntime.Runtime

	// sizing
	width, height       int
	subwidth, subheight int

	// keymap & help
	keymap rootKeymap
	help   help.Model

	// curr view
	curr      tea.Model
	currIdx   int
	currName  string
	currTitle string
	currStyle lipgloss.Style

	// views
	// main mainModel
	// TODO, map this?
	dash *dashModel
	list *listModel
	info *infoModel
	chat *chatModel

	// other fields
	msg string
	err error

	// session bookkeeping
	currSid       string
	currSessTitle string
	session       session.Session
	asession      *common.Session

	// preset / options bookkeeping
	currPreset string
	currAgent  string
	currModel  string
	currEnv    string
	currDir    string
}

var views = []string{"list", "chat", "info"}

func InitialModel(R *runtime.Runtime, AR *aruntime.Runtime, view string) (*Model, error) {
	w := 50
	h := 25

	if view == "" {
		view = "list"
	}
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
		currName:  view,
		currStyle: cs,
		currIdx:   0,
		// sess:      sess,
	}

	switch view {
	case "list":
		vm := initialSessionsModel(m)
		m.list = vm
		m.curr = vm
	case "chat":
		vm := initialChatModel(m)
		m.chat = vm
		m.curr = vm
	default:
		return nil, fmt.Errorf("unknown view %q", view)
	}

	return m, nil
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
		// TODO, mode the input from list to this model, moves us towards a shared setup and possibility for bloomberg like UX
		switch {
		case key.Matches(msg, m.keymap.add):
			if (m.currName == "chat" && m.chat != nil && !m.chat.textarea.Focused()) || (m.currName == "list" && !m.list.input.Focused()) {
				// m.createSession()
				m.updateCurrName("chat")
				m.chat.textarea.Focus()
				return m, tea.Sequence(mainCmd)
			}

		case key.Matches(msg, m.keymap.help):
			if (m.currName == "chat" && m.chat != nil && !m.chat.textarea.Focused()) || (m.currName == "list" && !m.list.input.Focused()) {
				show := !m.help.ShowAll
				m.help.ShowAll = show
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
		} else if m.currSessTitle != "" {
			first = m.currSessTitle
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
		subhelp = m.help.View(m.list.keymap)
	case "info":
		subhelp = m.help.View(m.info.keymap)
	case "chat":
		subhelp = m.help.View(m.chat.keymap)
	}

	err := lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Render("ok")
	if m.err != nil {
		err = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render(m.err.Error())
	}

	header := lipgloss.JoinHorizontal(lipgloss.Top, subhelp, " | ", help, " | ", err, " | ", first)
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
			m.info = initialInfoModel(m)
		}
		m.curr = m.info
		m.info.renderTab()
	case "list":
		if m.list == nil {
			m.list = initialSessionsModel(m)
		}
		m.curr = m.list
		m.list.updateSessions()
		m.list.updateRows()
	case "chat":
		if m.chat == nil {
			m.chat = initialChatModel(m)
		}
		m.curr = m.chat
		m.chat.textarea.Focus()
		m.chat.updateChatStatus()
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
	m.session = nil
	m.asession = nil
	if m.chat != nil {
		m.chat.viewport.SetContent("")
	}
}

func (m *Model) currDims() {
	var hcnt, ccnt, icnt, scnt int
	htxt := m.help.View(m.keymap)
	hcnt = strings.Count(htxt, "\n") + 1

	if m.chat != nil && m.currName == "chat" {
		htxt = m.help.View(m.chat.keymap)
		ccnt = strings.Count(htxt, "\n") + 1
	}
	if m.info != nil && m.currName == "info" {
		htxt = m.help.View(m.info.keymap)
		icnt = strings.Count(htxt, "\n") + 1
	}
	if m.list != nil && m.currName == "list" {
		htxt = m.help.View(m.list.keymap)
		scnt = strings.Count(htxt, "\n") + 1
	}
	hcnt = max(hcnt, ccnt, icnt, scnt)

	vh := m.height - hcnt + 1
	m.subwidth = m.width
	m.subheight = vh
}

func (m *Model) createSession() error {
	m.clearSession()

	session, err := common.SessionCreate(m.R.Ctx, m.AR, common.CreatePayload{
		User:    consts.VEG_DEFAULT_USER,
		Agent:   m.currAgent,
		Model:   m.currModel,
		EnvName: m.currEnv,
		Environ: &environ.EnvironCreateOptions{
			SrcUri: m.currDir,
		},
	})
	if err != nil {
		m.err = err
		return err
	}
	return m.setSession(session)
}

func (m *Model) loadSession(sid string) error {
	session, err := common.SessionGet(m.R.Ctx, m.AR, consts.VEG_DEFAULT_USER, sid)
	if err != nil {
		m.err = err
		return err
	}
	return m.setSession(session)
}

func (m *Model) sendMessage(text string) error {
	p := &common.ChatPayload{
		User:  consts.VEG_DEFAULT_USER,
		Sid:   m.currSid,
		Agent: m.currAgent,
		Model: m.currModel,
		Text:  text,
	}
	s, err := common.SessionChat(m.R, m.AR, p)
	if err != nil {
		m.err = err
		return err
	}
	m.asession = s
	return nil
}

func (m *Model) delSession(sid string) error {
	err := common.SessionDel(m.R.Ctx, m.AR, consts.VEG_DEFAULT_USER, sid)
	if err != nil {
		m.err = err
		return err
	}
	m.clearSession()
	return nil
}

func (m *Model) setSession(session session.Session) error {
	m.session = session
	m.currSid = session.ID()

	m.currSessTitle = ""
	m.currAgent = ""
	m.currModel = ""
	m.currEnv = ""
	m.currDir = ""

	state := maps.Collect(session.State().All())
	if v, ok := state["title"]; ok {
		m.currSessTitle = v.(string)
	}
	if v, ok := state["agent"]; ok {
		m.currAgent = v.(string)
	}
	if v, ok := state["model"]; ok {
		m.currModel = v.(string)
	}
	if v, ok := state["envName"]; ok {
		m.currEnv = v.(string)
	}
	// this is probably wrong
	if v, ok := state["initEnv"]; ok {
		vm := v.(map[string]any)
		if vd, ok := vm["srcUri"]; ok {
			m.currDir = vd.(string)
		}
	}

	return nil
}

func (m *Model) updateRootTitle() {
	if m.session == nil {
		m.currTitle = "nil session"
		return
	}
	state := maps.Collect(m.session.State().All())
	numEvents := m.session.Events().Len()
	numState := len(state)
	title := m.currSessTitle
	if title == "" {
		title = m.currSid
	}
	m.currTitle = fmt.Sprintf("%s  events:%d  state:%d", title, numEvents, numState)
}
