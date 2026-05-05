import "list"

import "strings"

import "time"

engine: _

_b2g:    1024.0 * 1024.0 * 1024.0
self:    engine.localCache.entrySet
total:   self.DiskSpaceBytes / _b2g
entries: engine.localCache.entrySet.Self.EntriesList

_e: {
	_#e: _
	t:   time.Unix(0, _#e.CreatedTimeUnixNano)
	l:   time.Unix(0, _#e.MostRecentUseTimeUnixNano)
	s:   _#e.DiskSpaceBytes / _b2g
	d:   _#e.Description
}

e1: [for e in entries {_e & {_#e: e}}]
e2: [for e in entries if e.DiskSpaceBytes/_b2g > 0.01 {_e & {_#e: e}}]
e2t: list.Sum([for e in e2 {e.size}])

e3: [for e in entries if strings.Contains(e.Description, "copy upload .git") {_e & {_#e: e}}]
e3t: list.Sum([for e in e3 {e.size}])

e4: [for e in entries if strings.Contains(e.Description, "-o qa ./") || strings.Contains(e.Description, "src/utils/qtp_agent") {_e & {_#e: e}}]
e4t: list.Sum([for e in e4 {e.size}])
