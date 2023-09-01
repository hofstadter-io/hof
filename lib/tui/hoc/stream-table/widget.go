package streamtable

import (
	"sync"

	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/tview"
)

type StreamTableSource func(chan string) chan interface{}
type StreamTableFormatter func(interface{}) [][]*tview.TableCell

type StreamTable struct {
	sync.Mutex
	*tview.Table

	TableHeader   [][]*tview.TableCell
	DataSource    StreamTableSource
	DataFormatter StreamTableFormatter

	dataCommands chan string
	dataStreamer chan interface{}
	quitChan     chan int
}

func NewStreamTable(header [][]*tview.TableCell, source StreamTableSource, formatter StreamTableFormatter) *StreamTable {
	ST := &StreamTable{
		Table:         tview.NewTable(),
		TableHeader:   header,
		DataSource:    source,
		DataFormatter: formatter,
	}

	return ST
}

func (ST *StreamTable) StartStream() {
	ST.Lock()
	defer ST.Unlock()

	// already shown
	if ST.dataStreamer != nil {
		return
	}

	ST.quitChan = make(chan int)
	ST.dataCommands = make(chan string)
	ST.dataStreamer = ST.DataSource(ST.dataCommands)

	go func() {
		for {
			ST.Lock()
			ds := ST.dataStreamer
			ST.Unlock()
			select {
			case data := <-ds:
				ST.UpdateData(data)

			case <-ST.quitChan:
				ST.Lock()
				defer ST.Unlock()
				ST.dataCommands <- "quit"
				close(ST.dataCommands)
				close(ST.quitChan)
				ST.quitChan = nil
				ST.dataCommands = nil
				ST.dataStreamer = nil
				return
			}
		}
	}()
}
func (ST *StreamTable) StopStream() {
	ST.quitChan <- 1
}

func (ST *StreamTable) UpdateData(input interface{}) {

	ST.Lock()
	data := ST.DataFormatter(input)

	cells := [][]*tview.TableCell{}
	cells = append(cells, ST.TableHeader...)
	cells = append(cells, data...)

	ST.Unlock()

	ST.Table.SetCells(cells)

	tui.Draw()
}
