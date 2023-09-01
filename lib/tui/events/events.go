// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Copyright 2018 Tony Worm <verdverm@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package events

import (
	"path"
	"strings"
	"sync"
	"time"

	"github.com/codemodus/kace"
	"github.com/gdamore/tcell/v2"

	"github.com/hofstadter-io/hof/lib/tui/tview"
)

type Event struct {
	tcell.Event
	when time.Time

	Type string
	Path string

	From string
	To   string

	Data interface{}
}

func (E *Event) When() time.Time {
	if E.Event != nil {
		return E.Event.When()
	} else {
		return E.when
	}
}

type EventError struct {
	*tcell.EventError
}

type EventResize struct {
	*tcell.EventResize
}

type EventKey struct {
	*tcell.EventKey
	KeyStr string
}

type EventMouse struct {
	*tcell.EventMouse
	Press string
}

type EventInterrupt struct {
	*tcell.EventInterrupt
}

type EventCustom struct {
	*tcell.EventInterrupt
}

func (EBus *EventBus) NewSysEvtCh() chan Event {
	ec := make(chan Event, 0)
	EBus.systemEventChans = append(EBus.systemEventChans, ec)
	return ec
}

func (EBus *EventBus) SendCustomEvent(path string, data interface{}) {
	now := time.Now()
	c := &EventCustom{
		EventInterrupt: tcell.NewEventInterrupt(data),
	}
	e := Event{
		when: now,
		Type: "custom",
		From: "user",
		Path: path,
		Data: c,
	}

	EBus.customEventChan <- Event(e)
}

func (EBus *EventBus) hookEventsFromApp(app *tview.Application) {
	hook := func(e tcell.Event) tcell.Event {
		for _, c := range EBus.systemEventChans {
			func(ch chan Event) {
				ch <- handleEvents(e)
			}(c)
		}
		return e
	}
	app.SetEventCapture(hook)
}

func handleEvents(e tcell.Event) Event {
	ne := Event{Event: e, From: "/sys"}

	switch t := e.(type) {
	case *tcell.EventError:
		err := EventError{t}
		ne.Type = "error"
		ne.Path = "/sys/err"
		ne.Data = err

	case *tcell.EventInterrupt:
		i := EventInterrupt{t}
		ne.Type = "interrupt"
		ne.Path = "/sys/interrupt"
		ne.Data = i

	case *tcell.EventResize:
		r := EventResize{t}
		ne.Type = "resize"
		ne.Path = "/sys/resize"
		ne.Data = r

	case *tcell.EventKey:
		k := eventKey(t)
		ne.Type = "keyboard"
		ne.Path = "/sys/key/" + k.KeyStr
		ne.Data = k

	case *tcell.EventMouse:
		m := eventMouse(t)
		ne.Type = "mouse"
		ne.Path = "/sys/mouse/" + m.Press
		ne.Data = m
	}
	return ne
}

func eventKey(e *tcell.EventKey) (tk EventKey) {
	tk.EventKey = e

	mods := eventMods(e.Modifiers())

	if e.Key() == tcell.KeyRune {
		k := string(e.Rune())
		tk.KeyStr = mods + k
		return tk
	}

	key := tcell.KeyNames[e.Key()]
	if strings.HasPrefix(key, "Ctrl-") {
		key = strings.TrimPrefix(key, "Ctrl-")
		key = kace.Kebab(key)
		if key == "space" {
			key = "<space>"
		}
	} else {
		key = "<" + kace.Kebab(key) + ">"
	}

	tk.KeyStr = mods + key

	return tk

}

func eventMouse(e *tcell.EventMouse) (te EventMouse) {
	te.EventMouse = e

	mods := eventMods(e.Modifiers())

	btn := ""
	switch B := e.Buttons(); B {
	case B & tcell.Button1:
		btn = "<left>"
	case B & tcell.Button2:
		btn = "<middle>"
	case B & tcell.Button3:
		btn = "<right>"

	case B & tcell.Button4:
		btn = "<button-4>"
	case B & tcell.Button5:
		btn = "<button-5>"

	case B & tcell.WheelUp:
		btn = "<wheel-up>"
	case B & tcell.WheelDown:
		btn = "<wheel-down>"
	case B & tcell.WheelLeft:
		btn = "<wheel-left>"
	case B & tcell.WheelRight:
		btn = "<wheel-right>"

	default:
		btn = "<unknown>"
	}

	te.Press = mods + btn

	return te
}

func eventMods(mods tcell.ModMask) string {
	m := ""
	if mods&tcell.ModCtrl == tcell.ModCtrl {
		m += "C-"
	}
	if mods&tcell.ModShift == tcell.ModShift {
		m += "S-"
	}
	if mods&tcell.ModAlt == tcell.ModAlt {
		m += "A-"
	}
	if mods&tcell.ModMeta == tcell.ModMeta {
		m += "M-"
	}
	return m
}

type EventStream struct {
	sync.RWMutex
	srcMap      map[string]chan Event
	stream      chan Event
	wg          sync.WaitGroup
	sigStopLoop chan Event
	Handlers    map[string]func(Event)
	hook        func(Event)
}

func NewEventStream() *EventStream {
	return &EventStream{
		srcMap:      make(map[string]chan Event),
		stream:      make(chan Event, 256),
		Handlers:    make(map[string]func(Event)),
		sigStopLoop: make(chan Event),
	}
}

func (es *EventStream) Init() {
	es.Merge("internal", es.sigStopLoop)
	go func() {
		es.wg.Wait()
		close(es.stream)
	}()
}

func (es *EventStream) Merge(name string, ec chan Event) {
	es.Lock()
	defer es.Unlock()

	es.wg.Add(1)
	es.srcMap[name] = ec

	go func(a chan Event) {
		for n := range a {
			n.From = name
			es.stream <- n
		}
		es.wg.Done()
	}(ec)
}

func (es *EventStream) Handle(path string, handler func(Event)) {
	es.Handlers[cleanPath(path)] = handler
}

func (es *EventStream) RemoveHandle(path string) {
	delete(es.Handlers, cleanPath(path))
}

// Remove all existing defined Handlers from the map
func (es *EventStream) ResetHandlers() {
	for Path := range es.Handlers {
		delete(es.Handlers, Path)
	}
	return
}

func (es *EventStream) match(path string) string {
	return findMatch(es.Handlers, path)
}

func (es *EventStream) Hook(f func(Event)) {
	es.hook = f
}

func (es *EventStream) Loop() {
	for e := range es.stream {
		switch e.Path {
		case "/sig/stoploop":
			return
		}
		func(a Event) {
			es.RLock()
			defer es.RUnlock()
			if pattern := es.match(a.Path); pattern != "" {
				es.Handlers[pattern](a)
			}
		}(e)

		if es.hook != nil {
			es.hook(e)
		}
	}

}

func (es *EventStream) StopLoop() {
	go func() {
		e := Event{
			Path: "/sig/stoploop",
		}
		es.sigStopLoop <- e
	}()
}

func findMatch(mux map[string]func(Event), path string) string {
	n := -1
	pattern := ""
	for m := range mux {
		if !isPathMatch(m, path) {
			continue
		}
		if len(m) > n {
			pattern = m
			n = len(m)
		}
	}
	return pattern

}

func cleanPath(p string) string {
	if p == "" {
		return "/"
	}
	if p[0] != '/' {
		p = "/" + p
	}
	return path.Clean(p)
}

func isPathMatch(pattern, path string) bool {
	if len(pattern) == 0 {
		return false
	}
	n := len(pattern)
	return len(path) >= n && path[0:n] == pattern
}
