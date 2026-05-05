package tui

import (
	"fmt"
	"maps"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/hofstadter-io/hof/lib/agent/config"
)

type chatModel struct {
	root   *Model
	keymap chatKeymap

	// model specific
	chatStyle  lipgloss.Style
	viewport   viewport.Model
	textarea   textarea.Model
	userStyle  lipgloss.Style
	agentStyle lipgloss.Style
	funcStyle  lipgloss.Style
	markglam   *glamour.TermRenderer
	chatStatus string
}

const CHAT_HEIGHT_TRIM_LENGTH = 6
const CHAT_INPUT_HEIGHT = 8

func initialChatModel(root *Model) *chatModel {
	width, height := root.subwidth, root.subheight
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 8192

	ta.SetWidth(width)
	ta.SetHeight(CHAT_INPUT_HEIGHT)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false
	vh := height - ta.Height() - CHAT_HEIGHT_TRIM_LENGTH

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
		textarea:   ta,
		viewport:   vp,
		userStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("4")),
		agentStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("2")),
		funcStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("3")),
		markglam:   glam,
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

	m.viewport.Height = m.root.subheight - m.textarea.Height() - CHAT_HEIGHT_TRIM_LENGTH
	m.viewport.Width = m.root.subwidth
	m.textarea.SetWidth(m.root.subwidth)

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

	if m.root.session != nil {
		m.root.updateRootTitle()
		m.updateMessagesFromEvents()
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// m.root.msg = msg.String()
		switch {
		case key.Matches(msg, m.keymap.info):
			if !m.textarea.Focused() {
				m.root.updateCurrName("info")
			}

		case key.Matches(msg, m.keymap.back):
			if m.textarea.Focused() {
				m.textarea.Blur()
			} else {
				m.root.clearSession()
				m.root.updateCurrName("list")
			}
		case key.Matches(msg, m.keymap.send):
			if m.textarea.Focused() {
				content := m.textarea.Value()
				lines := strings.Split(content, "\n")

				// look at first char and make a call
				fchar := lines[0][:1]
				switch fchar {

				// handle command parts
				case "$": // internal command
					for l, line := range lines {
						line = strings.TrimSpace(line)
						if line == "" {
							continue
						}
						parts := strings.Fields(line)
						fpart := parts[0][1:]
						switch fpart {
						case "agent":
							if len(parts) == 1 {
								m.root.currAgent = ""
							} else {
								m.root.currAgent = parts[1]
							}
						case "model":
							if len(parts) == 1 {
								m.root.currModel = ""
							} else {
								m.root.currModel = parts[1]
							}
						case "env":
							if len(parts) == 1 {
								m.root.currEnv = ""
							} else {
								m.root.currEnv = parts[1]
							}
						case "dir":
							if len(parts) == 1 {
								m.root.currDir = ""
							} else {
								pd := parts[1]

								if strings.HasPrefix(pd, ".") || strings.HasPrefix(pd, "/") {
									d, err := filepath.Abs(pd)
									if err != nil {
										m.root.err = err
										return m, nil
									}
									m.root.currDir = "file://" + d
								}
								if strings.HasPrefix(pd, "https://") {
									m.root.currDir = pd
								}
							}

						case "preset":
							preset := "default"
							if len(parts) == 2 {
								preset = parts[1]
							}
							if preset == "clear" {
								m.root.currAgent = ""
								m.root.currModel = ""
								m.root.currEnv = ""
								m.root.currDir = ""
							}

							var v cue.Value
							for _, ag := range m.root.R.Agentics {
								if ag.Hof.Agentic.Kind == "preset" && ag.Hof.Metadata.Name == preset {
									v = ag.Value
									break
								}
							}
							if !v.Exists() {
								m.root.err = fmt.Errorf("unknown preset %q", preset)
								return m, nil
							}
							var p config.Preset
							v.Decode(&p)
							m.root.currAgent = p.Agent
							m.root.currModel = p.Model
							m.root.currEnv = p.Env
							if strings.HasPrefix(p.Dir, ".") || strings.HasPrefix(p.Dir, "/") {
								d, err := filepath.Abs(p.Dir)
								if err != nil {
									m.root.err = err
									return m, nil
								}
								m.root.currDir = "file://" + d
							}
							if strings.HasPrefix(p.Dir, "https://") {
								m.root.currDir = p.Dir
							}

						// start new session
						case "chat":
							if len(parts) == 2 {
								m.root.currTitle = parts[1]
							}
							m.root.createSession()

						// modify state
						case "state":

						default:
							// user error
							m.root.err = fmt.Errorf("unknown command %q on line %d", fpart, l)
							return m, nil

						}
					}
					m.root.msg = "success!"
					m.updateChatStatus()
					m.textarea.Reset()
					return m, nil

				case "!": // run command, !! to add a layer

				// skills & such
				case "/":

				// chat message
				default:
					if m.root.session != nil {

						// sending first message
						m.root.msg = lines[0]
						err := m.root.sendMessage(content)
						if err != nil {
							// fmt.Println("root.sendMessage.error", err)
							m.root.err = err
							return m, nil
						}

						// try to get the UI to refresh
						m.textarea.Blur()

						m.textarea.Reset()

						// call loadSession on new events so we get all the info
						go func() {
							sess := m.root.asession
							// every time we get an event...
							for _ = range sess.EventChan {
								// fmt.Println("event:", e.Author, e.ID)
								// load the full lasest (being lazy, but also don't have to deal with state deltas and updates)
								// we should send the latest state back or make a func/api for returning session details w/o event list (between list & get today)
								err := m.root.loadSession(m.root.currSid)
								if err != nil {
									// fmt.Println("sess.ChatLoop.error:", err)
									m.root.err = err
									continue
								}

								m.updateMessagesFromEvents()
								m.updateChatStatus()

								// look for any errors
								select {
								case err := <-sess.ErrorChan:
									// fmt.Println("sess.ErrorChat:", err)
									m.root.err = err
								default:
								}
								// fmt.Println("loop around")
							}
						}()
					} else {
						m.root.err = fmt.Errorf("no session, use $chat [title]")
					}

				}

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
		m.root.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m *chatModel) View() string {
	if m.textarea.Focused() {
		m.textarea.Prompt = m.agentStyle.Render("┃ ")
	} else {
		m.textarea.Prompt = "┃ "
	}

	s := fmt.Sprintf(
		"%s\n\n%s\n%s\n%s",
		m.viewport.View(),
		fmt.Sprintln(m.root.msg),
		m.chatStatus,
		m.textarea.View(),
	)

	ws := m.chatStyle.Width((m.root.subwidth - windowStyle.GetHorizontalFrameSize()))
	return ws.Render(s)
}

func (m *chatModel) refresh() {
	m.root.loadSession(m.root.currSid)
	m.root.updateRootTitle()
	m.updateChatStatus()
	m.updateMessagesFromEvents()
	m.viewport.GotoBottom()
}

func (m *chatModel) updateChatStatus() {
	// if m.root.session == nil {
	// 	m.chatStatus = "nil session"
	// 	return
	// }

	var title, agent, model, envName, dirName string
	state := make(map[string]any)

	if m.root.session != nil {
		id := m.root.session.ID()
		state = maps.Collect(m.root.session.State().All())
		if t, ok := state["title"]; ok {
			title = t.(string)
		} else {
			title = id
		}
	} else {
		title = "no session"
	}

	// prefer local, then state
	if m.root.currAgent != "" {
		agent = m.root.currAgent
	} else {
		if a, ok := state["agent"]; ok {
			agent = a.(string)
		}
	}

	if m.root.currModel != "" {
		model = m.root.currModel
	} else {
		if m, ok := state["model"]; ok {
			model = m.(string)
		}
	}

	if m.root.currEnv != "" {
		envName = m.root.currEnv
	} else {
		if e, ok := state["envName"]; ok {
			envName = e.(string)
		}
	}

	if m.root.currDir != "" {
		dirName = m.root.currDir
	} else {
		// this is probably wrong
		if v, ok := state["initEnv"]; ok {
			vm := v.(map[string]any)
			if vd, ok := vm["srcUri"]; ok {
				dirName = vd.(string)
			}
		}
	}

	title = lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Render(title)
	agent = lipgloss.NewStyle().Foreground(lipgloss.Color("4")).Render(agent)
	model = lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Render(model)
	envName = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Render(envName)
	dirName = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Render(dirName)

	m.chatStatus = lipgloss.JoinHorizontal(lipgloss.Top, " ", title, " | ", agent, " | ", model, " | ", envName, " | ", dirName) + "\n"
}

func (m *chatModel) updateMessagesFromEvents() {
	messages := renderMessages(m.root.subwidth, m.root.session)
	content := strings.Join(messages, "\n") + "\n\n\n\n"
	m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(content))
}
