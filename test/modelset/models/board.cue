package models

import (
	"github.com/hofstadter-io/dma/schema"
)

#Board: schema.#Model @modelsets(kanban)
#Board: {
	Name: "board" @db(sql)
}
