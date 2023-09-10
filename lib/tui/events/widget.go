// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package events

import (
	"fmt"
	"sync"

	"github.com/hofstadter-io/hof/lib/tui/tview"
)

// event mixins
type WgtMgr map[string]WgtInfo

type WgtInfo struct {
	Handlers map[string]func(Event)
	WgtRef   tview.Primitive
	Id       string
	Lock     *sync.RWMutex
}

func NewWgtInfo(wgt tview.Primitive) WgtInfo {
	return WgtInfo{
		Handlers: make(map[string]func(Event)),
		WgtRef:   wgt,
		Id:       wgt.Id(),
		Lock:     new(sync.RWMutex),
	}
}

func NewWgtMgr() WgtMgr {
	wm := WgtMgr(make(map[string]WgtInfo))
	return wm

}

func (wm WgtMgr) AddWgt(wgt tview.Primitive) {
	wm[wgt.Id()] = NewWgtInfo(wgt)
}

func (wm WgtMgr) RmWgt(wgt tview.Primitive) {
	wm.RmWgtById(wgt.Id())
}

func (wm WgtMgr) RmWgtById(id string) {
	delete(wm, id)
}

func (wm WgtMgr) AddWgtHandler(id, path string, h func(Event)) {
	if w, ok := wm[id]; ok {
		w.Handlers[path] = h
	}
}

func (wm WgtMgr) RmWgtHandler(id, path string) {
	if w, ok := wm[id]; ok {
		delete(w.Handlers, path)
	}
}

func (wm WgtMgr) ClearWgtHandlers(id string) {
	if w, ok := wm[id]; ok {
		w.Handlers = make(map[string]func(Event))
	}
}

var counter struct {
	sync.RWMutex
	count int
}

func GenId() string {
	counter.Lock()
	defer counter.Unlock()

	counter.count += 1
	return fmt.Sprintf("%d", counter.count)
}

func (wm WgtMgr) WgtHandlersHook() func(Event) {
	return func(e Event) {
		for _, v := range wm {
			v.Lock.RLock()
			defer v.Lock.RUnlock()
			if k := findMatch(v.Handlers, e.Path); k != "" {
				// v.WgtRef.Lock()
				v.Handlers[k](e)
				// v.WgtRef.Unlock()
			}
		}
	}
}
