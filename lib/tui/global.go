package tui

import (
	"github.com/hofstadter-io/hof/lib/tui/app"
	"github.com/hofstadter-io/hof/lib/tui/events"
	"github.com/hofstadter-io/hof/lib/tui/tview"
)

// this is where we handle a few global things
// eventually this needs to be refactored...
// but for now, we only support one app, which needs to be set
// but at least you can have a few and swap them out

var globalApp *app.App

func GetApp()(a *app.App) {
	return globalApp
}

func SetApp(a *app.App) {
	globalApp = a
}

func Draw() {
	go globalApp.Draw()
}

func Clear() {
	globalApp.Clear()
}

func Stop() {
	globalApp.Stop()
}

func GetRootView() tview.Primitive {
	return globalApp.GetRootView()
}

func SetRootView(v tview.Primitive) {
	globalApp.SetRootView(v)
}

func GetFocus() (p tview.Primitive) {
	if globalApp == nil {
		return nil
	}
	return globalApp.GetFocus()
}

func SetFocus(p tview.Primitive) {
	//appLock.Lock()
	//defer appLock.Unlock()

	if globalApp == nil {
		// really shouldn't get here, but the event stream is still running
		return
	}

	// go app.Screen().HideCursor()
	globalApp.SetFocus(p)
	Draw()
}

func Unfocus() {
	//appLock.Lock()
	//defer appLock.Unlock()

	if globalApp == nil {
		// really shouldn't get here, but the event stream is still running
		return
	}

	// go app.Screen().HideCursor()
	globalApp.SetFocus(globalApp.GetRootView())
	Draw()
}

func QueueUpdate(f func()) {
	globalApp.QueueUpdate(f)
}

func QueueUpdateDraw(f func()) {
	globalApp.QueueUpdateDraw(f)
}

func SendCustomEvent(path string, data any) {
	globalApp.EventBus.SendCustomEvent(path, data)
}

func Log(level string, data any) {
	if level == "crit" {
		globalApp.EventBus.SendCustomEvent("/user/error", data)
	}
	globalApp.EventBus.SendCustomEvent("/console/" + level, data)
}

func Tell(level string, data any) {
	globalApp.EventBus.SendCustomEvent("/user/" + level, data)
}

func AddGlobalHandler(path string, handler func(events.Event)) {
	globalApp.EventBus.AddGlobalHandler(path, handler)
}

func RemoveGlobalHandler(path string) {
	globalApp.EventBus.RemoveGlobalHandler(path)
}

func AddWidgetHandler(widget tview.Primitive, path string, handler func(events.Event)) {
	globalApp.EventBus.AddWidgetHandler(widget, path, handler)
}

func RemoveWidgetHandler(widget tview.Primitive, path string) {
	globalApp.EventBus.RemoveWidgetHandler(widget, path)
}

func ClearWidgetHandlers(widget tview.Primitive) {
	globalApp.EventBus.ClearWidgetHandlers(widget)
}
