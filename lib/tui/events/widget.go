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
type WgtMgr struct {
	sync.Map // WgtInfo
}

type WgtInfo struct {
	handlers sync.Map // func(Event)
	wgtRef   tview.Primitive
	id       string
}

func NewWgtInfo(wgt tview.Primitive) *WgtInfo {
	return &WgtInfo{
		wgtRef:   wgt,
		id:       wgt.Id(),
	}
}

func (wm *WgtMgr) AddWgt(wgt tview.Primitive) {
	wm.Store(wgt.Id(), NewWgtInfo(wgt))
}

func (wm *WgtMgr) RmWgt(wgt tview.Primitive) {
	wm.RmWgtById(wgt.Id())
}

func (wm *WgtMgr) RmWgtById(id string) {
	wm.Delete(id)
}

func (wm *WgtMgr) AddWgtHandler(id, path string, h func(Event)) {
	if w, ok := wm.Load(id); ok {
		W := w.(*WgtInfo)
		W.handlers.Store(path, h)
	}
}

func (wm *WgtMgr) RmWgtHandler(id, path string) {
	if w, ok := wm.Load(id); ok {
		W := w.(*WgtInfo)
		W.handlers.Delete(path)
	}
}

func (wm *WgtMgr) ClearWgtHandlers(id string) {
	if w, ok := wm.Load(id); ok {
		W := w.(*WgtInfo)
		W.handlers = sync.Map{}
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

func (wm *WgtMgr) WgtHandlersHook() func(Event) {
	return func(e Event) {
		wm.Range(func (wk, wv any) bool {
			WV := wv.(*WgtInfo)
			var H func(Event)

			n := -1
			WV.handlers.Range(func (hk, hv any) bool {
				m := hk.(string)

				if !isPathMatch(m, e.Path) {
					return true
				}
				// looking for longest match
				if len(m) > n {
					H = hv.(func(Event))
				}

				return true
			})

			// if we had a match, lets call the event handler
			if H != nil {
				H(e)
			}

			return true
		})
	}
}
