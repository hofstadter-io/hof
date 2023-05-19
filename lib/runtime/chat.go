package runtime

import (
	"fmt"

	"github.com/hofstadter-io/hof/lib/chat"
)

type ChatEnricher func(*Runtime, *chat.Chat) error

func (R *Runtime) EnrichChats(chats []string, enrich ChatEnricher) error {
	if R.Flags.Verbosity > 1 {
		fmt.Println("Runtime.Chat: ", chats)
		for _, node := range R.Nodes {
			node.Print()
		}
	}

	// Find only the datamodel nodes
	// TODO, dedup any references
	cs := []*chat.Chat{}
	for _, node := range R.Nodes {
		// check for DM root
		if node.Hof.Chat.Root {

			cs = append(cs, &chat.Chat{Node: node})
		}
	}

	R.Chats = cs

	for _, c := range R.Chats {
		err := enrich(R, c)
		if err != nil {
			return err
		}
	}


	return nil
}
