package cmdbox

import (
	"fmt"
	"strings"
	"sync"

	"github.com/gdamore/tcell/v2"

	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/events"
	"github.com/hofstadter-io/hof/lib/tui/tview"
)

const emptyMsg = "press 'Ctrl-<space>' to enter a command or '/path/to/something' to navigate"

type Command interface {
	CommandName() string
	CommandUsage() string
	CommandHelp() string

	CommandCallback(context map[string]any)
}

type DefaultCommand struct {
	Name  string
	Usage string
	Help  string

	Callback func(context map[string]interface{})
}

func (DC *DefaultCommand) CommandName() string {
	return DC.Name
}

func (DC *DefaultCommand) CommandHelp() string {
	return DC.Help
}

func (DC *DefaultCommand) CommandUsage() string {
	return DC.Usage
}

func (DC *DefaultCommand) CommandCallback(context map[string]interface{}) {
	DC.Callback(context)
}

type CmdBoxWidget struct {
	*tview.InputField

	sync.Mutex

	commands map[string]Command

	curr    string   // current input (potentially partial)
	hIdx    int      // where we are in history
	history []string // command history

	lastFocus tview.Primitive

	getLastCmd func() string
	setLastCmd func(string)
}

func New(getLastCmd func() string, setLastCmd func(string)) *CmdBoxWidget {
	cb := &CmdBoxWidget{
		InputField: tview.NewInputField(),
		commands:   make(map[string]Command),
		history:    []string{},
		getLastCmd: getLastCmd,
		setLastCmd: setLastCmd,
	}

	cb.InputField.
		SetFieldBackgroundColor(tview.Styles.PrimitiveBackgroundColor).
		SetLabelText(" ")

	return cb
}

func (CB *CmdBoxWidget) Id() string {
	return CB.InputField.Id()
}

func (CB *CmdBoxWidget) AddCommandCallback(command string, callback func(map[string]interface{})) Command {
	CB.Lock()
	defer CB.Unlock()
	c := &DefaultCommand{
		Name:     command,
		Usage:    command,
		Help:     "no help for " + command,
		Callback: callback,
	}
	CB.commands[c.CommandName()] = c
	return c
}

func (CB *CmdBoxWidget) AddCommand(command Command) {
	CB.Lock()
	defer CB.Unlock()
	// go tui.SendCustomEvent("/console/info", "adding command: "+command.CommandName())
	CB.commands[command.CommandName()] = command
}

func (CB *CmdBoxWidget) RemoveCommand(command Command) {
	delete(CB.commands, command.CommandName())
}

func (CB *CmdBoxWidget) Mount(context map[string]interface{}) error {
	tui.AddWidgetHandler(CB, "/sys/key/C-<space>", func(e events.Event) {
		CB.Lock()
		CB.curr = ""
		CB.hIdx = len(CB.history)
		CB.Unlock()

		CB.SetText("")

		CB.lastFocus = tui.GetFocus()
		tui.SetFocus(CB.InputField)
	})

	CB.SetFinishedFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			input := CB.GetText()
			input = strings.TrimSpace(input)
			if input != "" {
				flds := strings.Fields(input)
			
				if CB.lastFocus != nil {
					tui.SetFocus(CB.lastFocus)
					CB.lastFocus = nil
				} else {
					tui.Unfocus()
				}

				CB.Submit(flds[0], flds[1:])
				CB.SetText("")
			}
		case tcell.KeyEscape:
			CB.SetText("")
			if CB.lastFocus != nil {
				tui.SetFocus(CB.lastFocus)
				CB.lastFocus = nil
			} else {
				tui.Unfocus()
			}

		// reserved for autocomplete
		case tcell.KeyTab:
		case tcell.KeyBacktab:

		default:
			go tui.SendCustomEvent("/console/warn", fmt.Sprintf("cmdbox (fin-???-key): %v", key))

		}

	})

	return nil
}
func (CB *CmdBoxWidget) Unmount() error {
	tui.RemoveWidgetHandler(CB, "/sys/key/C-<space>")

	return nil
}

func (CB *CmdBoxWidget) Submit(command string, args []string) {
	if len(command) == 0 {
		return
	}

	CB.Lock()
	if len(args) == 0 {
		CB.history = append(CB.history, command)
	} else {
		CB.history = append(CB.history, command+" "+strings.Join(args, " "))
	}
	CB.Unlock()

	command = strings.ToLower(command)
	// exact route?
	if command[:1] == "/" {
		go tui.SendCustomEvent("/router/dispatch", command)
		return
	}

	var cmd Command

	if command[:1] == ":" {
		// global command
		command = command[1:]
	}	else {
		// local command
		// move command to head of args
		args = append([]string{command}, args...)
		// set the actual command to the last one
		command = CB.getLastCmd()
	}

	tui.Log("debug", fmt.Sprintf("CB.Submit %q %v", command, args))

	CB.Lock()
	var ok bool
	cmd, ok = CB.commands[command]
	CB.Unlock()
	if !ok {
		tui.Unfocus()
		// render for the user
		go tui.SendCustomEvent("/user/error", fmt.Sprintf("unknown command %q", command))
		// log to console
		go tui.SendCustomEvent("/console/warn", fmt.Sprintf("unknown command %q", command))
		return
	}

	ctx := map[string]any{ "args": args }
	CB.setLastCmd(command)
	go cmd.CommandCallback(ctx)
}

// InputHandler returns the handler for this primitive.
func (CB *CmdBoxWidget) InputHandler() func(*tcell.EventKey, func(tview.Primitive)) {
	return CB.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		handle := CB.InputField.InputHandler()

		dist := 1

		// Process key evt.
		switch key := event.Key(); key {

		// Upwards, back in history
		case tcell.KeyHome:
			dist = len(CB.history)
			fallthrough
		case tcell.KeyPgUp:
			dist += 4
			fallthrough
		case tcell.KeyUp: // Regular character.
			if CB.hIdx == len(CB.history) {
				CB.curr = CB.GetText()
			}
			CB.hIdx -= dist
			if CB.hIdx < 0 {
				CB.hIdx = 0
			}
			if len(CB.history) > 0 {
				CB.SetText(CB.history[CB.hIdx])
			}

		// Downwards, more recent in history
		case tcell.KeyEnd:
			dist = len(CB.history)
			fallthrough
		case tcell.KeyPgDn:
			dist += 4
			fallthrough
		case tcell.KeyDown:
			CB.hIdx += dist
			if CB.hIdx > len(CB.history) {
				CB.hIdx = len(CB.history)
			}

			if CB.hIdx == len(CB.history) {
				CB.SetText(CB.curr)
				return
			}
			if len(CB.history) > 0 {
				CB.SetText(CB.history[CB.hIdx])
			}

		// Default is to pass through to InputField handler
		default:
			CB.hIdx = len(CB.history)
			handle(event, setFocus)

		}
	})
}
