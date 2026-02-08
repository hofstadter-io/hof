package tui

import (
	"fmt"
	"maps"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hofstadter-io/hof/lib/agent/runtime/handlers/common"
	"github.com/hofstadter-io/hof/lib/consts"
	"google.golang.org/adk/session"
)

const SESSION_LIST_ID_POS = 4

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type listModel struct {
	root   *Model
	keymap listKeymap

	// model specific
	mode     string
	table    table.Model
	input    textinput.Model
	sorts    []string
	sessions []session.Session
	info     *infoModel
}

func initialSessionsModel(root *Model) *listModel {
	width, height := root.subwidth, root.subheight

	h := height - 9
	mw := width - 3

	t := table.New(
		// table.WithColumns(columns),
		// table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(h),
		table.WithWidth(mw),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	i := textinput.New()

	m := &listModel{
		root:   root,
		keymap: sessionsKeymapDefaults,
		table:  t,
		input:  i,
		sorts:  []string{"update"},
	}

	m.updateSessions()
	m.updateColumns()
	m.updateRows()

	return m
}

func (m *listModel) Init() tea.Cmd { return nil }

func (m *listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// calc sub dims
	h := m.root.subheight - 3
	mw := m.root.subwidth

	// update dims
	m.table.SetHeight(h)
	m.table.SetWidth(mw)

	m.updateColumns()

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// if the input is focused*
		if m.input.Focused() {
			switch {
			// *no-ops, handled in the next switch statement
			case key.Matches(msg, m.keymap.back):
			case key.Matches(msg, m.keymap.load):

			// handle inputs
			default:
				m.input, cmd = m.input.Update(msg)
				return m, cmd
			}
		}

		// if the input is not focused*
		switch {
		case key.Matches(msg, m.keymap.info):
			if len(m.table.Rows()) > 0 {
				m.mode = "info"
				sid := m.table.SelectedRow()[SESSION_LIST_ID_POS]
				m.root.currSid = sid
				m.root.loadSession(sid)
				m.root.updateCurrName("info")
				m.root.updateRootTitle()
			}

		case key.Matches(msg, m.keymap.del):
			if len(m.table.Rows()) > 0 {
				sid := m.table.SelectedRow()[SESSION_LIST_ID_POS]
				m.root.delSession(sid)
				err := m.updateSessions()
				if err != nil {
					m.root.err = err
				}
				m.updateRows()
			}

		case key.Matches(msg, m.keymap.sort):
			m.mode = "sort"
			m.input.Focus()

		case key.Matches(msg, m.keymap.refresh):
			err := m.updateSessions()
			if err != nil {
				m.root.err = err
			}
			m.updateRows()

		case key.Matches(msg, m.keymap.back):
			switch m.mode {
			case "sort":
				m.input.Reset()
				m.input.Blur()
				m.mode = ""
			case "info":
				m.mode = ""
				m.root.clearSession()
				m.root.updateCurrName("list")
				cmd = m.root.runTick()
			}
		case key.Matches(msg, m.keymap.load):
			switch m.mode {
			case "sort":
				in := m.input.Value()
				ins := strings.Fields(in)
				m.sorts = ins
				if len(m.sorts) == 0 {
					m.sorts = []string{"update"}
				}
				m.input.Reset()
				m.input.Blur()
				m.mode = ""
				m.updateRows()

			default:
				if len(m.table.Rows()) > 0 {
					m.root.clearSession()
					sid := m.table.SelectedRow()[SESSION_LIST_ID_POS]
					m.root.currSid = sid
					m.root.updateCurrName("chat")
					m.root.chat.refresh()
					cmd = m.root.runTick()
				}
			}
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m *listModel) View() string {
	t := baseStyle.Render(m.table.View()) + "\n"

	switch m.mode {
	case "sort":
		return fmt.Sprintf("%s\n%s", m.input.View(), t)
	default:
		return t
	}
}

func (m *listModel) updateColumns() {
	cw := (m.root.subwidth - 9) / 3
	columns := []table.Column{
		{Title: "Pos", Width: 3},
		{Title: "Title", Width: cw},
		{Title: "State", Width: 6},
		{Title: "Last Update", Width: cw},
		{Title: "ID", Width: cw},
	}
	m.table.SetColumns(columns)
}

func (m *listModel) updateRows() {
	type datum struct {
		title  string
		id     string
		state  int
		create time.Time
		update time.Time
	}
	data := make([]datum, 0, len(m.sessions))
	for _, s := range m.sessions {
		id := s.ID()

		numState := len(maps.Collect(s.State().All()))
		tVal, err := s.State().Get("title")
		title := ""
		if err == nil {
			title = tVal.(string)
		}

		d := datum{
			title:  title,
			id:     id,
			state:  numState,
			update: s.LastUpdateTime(),
		}
		data = append(data, d)
	}

	sort.Slice(data, func(i, j int) bool {
		lhs, rhs := data[i], data[j]
		for _, s := range m.sorts {
			s = strings.ToLower(s)
			switch s {
			case "title", "name":
				return lhs.title < rhs.title
			case "-title", "-name":
				return lhs.title > rhs.title
			case "id":
				return lhs.id < rhs.id
			case "-id":
				return lhs.id < rhs.id
			case "state":
				return lhs.state < rhs.state
			case "-state":
				return lhs.state > rhs.state
			case "create":
				return lhs.create.After(rhs.create)
			case "-create":
				return lhs.create.Before(rhs.create)
			case "update":
				return lhs.update.After(rhs.update)
			case "-update":
				return lhs.update.Before(rhs.update)
			default:
				fmt.Printf("WARN: unknown sort field %q ignored", s)
				continue
			}
		}
		return false
	})

	rows := make([]table.Row, 0, len(data))
	for i, s := range data {
		name := s.title
		if name == "" {
			name = s.id
		}
		numState := fmt.Sprintf("%2d", s.state)
		lastTime := s.update.Local().Format("Mon, Jan 2, 2006 15:04")
		row := []string{fmt.Sprint(i), name, numState, lastTime, s.id}
		rows = append(rows, row)
	}

	m.table.SetRows(rows)

}

func (m *listModel) updateSessions() error {
	sessions, err := common.SessionList(m.root.R.Ctx, m.root.AR, consts.VEG_DEFAULT_USER)
	if err != nil {
		return err
	}

	m.sessions = sessions
	return nil
}
