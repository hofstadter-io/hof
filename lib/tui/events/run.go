package events

import (
	"sync"

	"github.com/hofstadter-io/hof/lib/tui/tview"
)

type EventBus struct {

	EventStream *EventStream
	WgtMgr WgtMgr
	wgtMgrMuxtx sync.Mutex

	Handler  func(e Event)

	systemEventChans []chan Event
	customEventChan  chan Event
	//= make(chan Event, 256)

}


func (EBus *EventBus) Init(app *tview.Application) error {
	EBus.systemEventChans = make([]chan Event, 0)
	EBus.customEventChan = make(chan Event, 256)
	go EBus.hookEventsFromApp(app)

	EBus.EventStream = NewEventStream()
	EBus.EventStream.Init()
	EBus.EventStream.Merge("tcell", EBus.NewSysEvtCh())
	EBus.EventStream.Merge("custom", EBus.customEventChan)

	EBus.WgtMgr = NewWgtMgr()
	EBus.EventStream.Hook(EBus.WgtMgr.WgtHandlersHook())

	return nil
}

func (EBus *EventBus) Start() error {
	EBus.EventStream.Loop()
	return nil
}

func (EBus *EventBus) Stop() error {
	EBus.EventStream.StopLoop()
	return nil
}

func (EBus *EventBus) Merge(name string, ec chan Event) {
	EBus.EventStream.Merge(name, ec)
}

func (EBus *EventBus) AddGlobalHandler(path string, handler func(Event)) {
	EBus.EventStream.Handle(path, handler)
}

func (EBus *EventBus) RemoveGlobalHandler(path string) {
	EBus.EventStream.RemoveHandle(path)
}

func (EBus *EventBus) ClearGlobalHandlers() {
	EBus.EventStream.ResetHandlers()
}

func (EBus *EventBus) AddWidgetHandler(wgt tview.Primitive, path string, handler func(Event)) {
	if _, ok := EBus.WgtMgr[wgt.Id()]; !ok {
		EBus.WgtMgr.AddWgt(wgt)
	}

	EBus.WgtMgr.AddWgtHandler(wgt.Id(), path, handler)
}

func (EBus *EventBus) RemoveWidgetHandler(wgt tview.Primitive, path string) {
	_, ok := EBus.WgtMgr[wgt.Id()]
	if !ok {
		return
	}

	EBus.WgtMgr.RmWgtHandler(wgt.Id(), path)
}

func (EBus *EventBus) ClearWidgetHandlers(wgt tview.Primitive) {
	_, ok := EBus.WgtMgr[wgt.Id()]
	if !ok {
		return
	}

	EBus.WgtMgr.ClearWgtHandlers(wgt.Id())
}
