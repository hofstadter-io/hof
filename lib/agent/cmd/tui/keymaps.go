package tui

import "github.com/charmbracelet/bubbles/key"

type rootKeymap struct {
	help, quit key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k rootKeymap) ShortHelp() []key.Binding {
	return []key.Binding{k.help, k.quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k rootKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.help, k.quit}, // first column
	}
}

var rootKeymapDefaults = rootKeymap{
	help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type listKeymap struct {
	add, back, info, load, refresh, sort key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k listKeymap) ShortHelp() []key.Binding {
	return []key.Binding{k.add, k.info, k.load}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k listKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.add, k.info, k.load}, // first column
		{k.sort, k.refresh, k.back},
	}
}

var sessionsKeymapDefaults = listKeymap{
	add: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new"),
	),
	back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	info: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "info"),
	),
	load: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "load"),
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
		key.WithHelp("←/h l/→", "nav"),
	),
	prev: key.NewBinding(
		key.WithKeys("left", "h"),
	),
	next: key.NewBinding(
		key.WithKeys("right", "l"),
	),
	chat: key.NewBinding(
		key.WithKeys("alt+enter"),
		key.WithHelp("alt+enter", "chat"),
	),
	back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
}

type chatKeymap struct {
	send, focus, back key.Binding
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
		{k.send, k.back}, // first column
	}
}

var chatKeymapDefaults = chatKeymap{
	send: key.NewBinding(
		key.WithKeys("alt+enter"),
		key.WithHelp("alt+enter", "send msg"),
	),
	focus: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "focus input"),
	),
	back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
}
