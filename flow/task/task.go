package task

import (
	"time"

	"cuelang.org/go/cue"
	cueflow "cuelang.org/go/tools/flow"
	"github.com/google/uuid"

	"github.com/hofstadter-io/hof/lib/hof"
)

type Task interface {
	IDer
	Eventer
	TimeEventer
}

type IDer interface {
	ID() string
	UUID() string
}

type Eventer interface {
	EmitEvent(key string, data interface{})
}

type TimeEventer interface {
	AddTimeEvent(key string)
}

type BaseTask struct {
	// IDer
	ID   string
	UUID uuid.UUID

	// cue bookkeeping
	CueTask *cueflow.Task
	Node    *hof.Node[any]
	Orig    cue.Value
	Start   cue.Value
	Final   cue.Value
	Error   error

	// stats & timing
	// should this be a list with names / times
	// timing
	// replace with open telemetry
	TimeEvents map[string]time.Time
}

func NewBaseTask(node *hof.Node[any]) *BaseTask {
	val := node.Value
	return &BaseTask{
		ID:         val.Path().String(),
		UUID:       uuid.New(),
		Node:       node,
		Orig:       val,
		TimeEvents: make(map[string]time.Time),
	}
}

func (T *BaseTask) AddTimeEvent(key string) {
	T.TimeEvents[key] = time.Now()
}
