package tui

import "github.com/charmbracelet/bubbles/key"

type rootKeymap struct {
	add, help, quit key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k rootKeymap) ShortHelp() []key.Binding {
	return []key.Binding{k.add, k.help, k.quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k rootKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.add, k.help, k.quit}, // first column
	}
}

var rootKeymapDefaults = rootKeymap{
	add: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new"),
	),
	help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type dashKeymap struct {
	nav, down, up, enter key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k dashKeymap) ShortHelp() []key.Binding {
	return []key.Binding{k.nav, k.enter}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k dashKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.nav, k.enter}, // first column

	}
}

var dashKeymapDefaults = dashKeymap{
	nav: key.NewBinding(
		key.WithKeys("doesnotexist"),
		key.WithHelp("↓/j ↑/k", "nav"),
	),
	down: key.NewBinding(
		key.WithKeys("down", "j"),
	),
	up: key.NewBinding(
		key.WithKeys("up", "k"),
	),
	enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
}

type listKeymap struct {
	del, back, info, load, refresh, sort key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k listKeymap) ShortHelp() []key.Binding {
	return []key.Binding{k.del, k.info, k.load}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k listKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.del, k.info, k.load}, // first column
		{k.sort, k.refresh, k.back},
	}
}

var sessionsKeymapDefaults = listKeymap{
	back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	del: key.NewBinding(
		key.WithKeys("D"),
		key.WithHelp("D", "del"),
	),
	info: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "info"),
	),
	load: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "chat"),
	),
	refresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refresh"),
	),
	sort: key.NewBinding(
		key.WithKeys(">"),
		key.WithHelp(">", "sort"),
	),
}

type infoKeymap struct {
	prev, next, nav, chat, back key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k infoKeymap) ShortHelp() []key.Binding {
	return []key.Binding{k.nav, k.chat, k.back}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k infoKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.nav, k.chat, k.back}, // first column

	}
}

var infoKeymapDefaults = infoKeymap{
	nav: key.NewBinding(
		key.WithKeys("doesnotexist"),
		key.WithHelp("←/h →/l", "nav"),
	),
	prev: key.NewBinding(
		key.WithKeys("left", "h"),
	),
	next: key.NewBinding(
		key.WithKeys("right", "l"),
	),
	chat: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "chat"),
	),
	back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
}

type chatKeymap struct {
	back, focus, info, send,
	agent, model, environ, open key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k chatKeymap) ShortHelp() []key.Binding {
	return []key.Binding{k.send, k.back}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k chatKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.send, k.back, k.info, k.focus},
		{k.agent, k.model, k.environ, k.open},
	}
}

var chatKeymapDefaults = chatKeymap{
	// navigation
	info: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "info"),
	),
	focus: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "focus input"),
	),
	back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),

	// send textarea as input
	send: key.NewBinding(
		key.WithKeys("alt+enter"),
		key.WithHelp("alt+enter", "send msg"),
	),

	// modals or something
	agent: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "agent"),
	),
	model: key.NewBinding(
		key.WithKeys("m"),
		key.WithHelp("m", "model"),
	),
	environ: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "environ"),
	),
	open: key.NewBinding(
		key.WithKeys("o"),
		key.WithHelp("o", "open"),
	),
}
