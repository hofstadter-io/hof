package runtime

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hofstadter-io/hof/lib/agent/config"
)

// Message is the "envelope" that all messages follow.
// The 'type' field tells us how to parse the 'payload'.
type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"` // Use RawMessage to delay payload parsing
}

// --- Client ---

// This should become more inclusive with optional features
// i.e. we should associate more with vscode connected to extension server
// i.e.2. we want to do cool things besides agentic with the virtualized FS backed by dagger
// planned integrations like kubernetes, dagger cache info, github/gerrit, build systems...
// (my ambitions for vscode extn development have exploded recently [as of writing this comment])

// Client is a wrapper for a single WebSocket connection (one VS Code window).
type Client struct {
	User  string
	State map[string]any // should this be persisted, do we even need it with user:... State? (same user on two clients, repo in different locations?)

	// when we have custom agents, or local to a session even? (b/c diff sess diff workdir)
	AgentDefs map[string]config.Agent

	// this really depends on the workspace / session
	// and should also be merged with (1) user global (2) builtin defaults
	// need a place for selecting which ones show up in the dropdown vs @mention [any]
	Agentic config.Config

	// we should perhaps store active sessions here
	// various information we'd like to share between agents (multiple vscode status/state)

	// other stuff needs to be persisted
	// 1. agent config (maybe we just store these in the state with user:...)
	// 2. session state/history (already done by SessionService, but needs improvements)

	conn *websocket.Conn

	send chan []byte // Buffered channel for outbound messages

	handleMessage func(*Client, *Message)
}

// readPump pumps messages from the WebSocket connection to the hub.
func (r *Runtime) readPump(c *Client) {
	defer func() {
		// On exit, unregister the client and close the connection
		r.unregister <- c
		c.conn.Close()
	}()

	// Set read limits, pong handlers, etc. (good practice)
	c.conn.SetReadLimit(100 * 1024 * 1024) // 100Mb (for passing files around)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(60 * time.Second)); return nil })

	var wg sync.WaitGroup

	for {
		// Read a message from the WebSocket
		_, jsonMessage, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("ERROR.unexpected-close: %v", err)
			}
			break // Exit loop on error
		}

		// Deserialize the JSON message envelope
		var msg Message
		if err := json.Unmarshal(jsonMessage, &msg); err != nil {
			log.Printf("Error unmarshaling message envelope: %v", err)
			continue // Keep processing other messages
		}

		// handleMessage is the main router for deserialized messages.
		// log.Printf("Received message type: %s", msg.Type)
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.handleMessage(c, &msg)
		}()
	}

	wg.Wait()
}

// writePump pumps messages from the hub to the WebSocket connection.
func (c *Client) writePump() {
	ticker := time.NewTicker(30 * time.Second) // Ping ticker
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			// The hub closed the channel.
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// prepare & write
			c.conn.SetWriteDeadline(time.Now().Add(60 * time.Second))
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// WARN, commented this out because it was causing JSON parse errors on the front end
			// the messages are not separated. Maybe we could parse them as JSONL (if everyone agrees to send json objs as a single line in their messages)
			// this generally seems redundent with the select statement above firing them off in rapid succession anyway
			// // Add queued chat messages to the current websocket message.
			// n := len(c.send)
			// for i := 0; i < n; i++ {
			// 	w.Write(<-c.send)
			// }

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			// Send ping
			c.conn.SetWriteDeadline(time.Now().Add(60 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) Mail(typ string, data any) {
	// Marshal the payload to JSON bytes
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling response: %v", err)
		return
	}

	// Create the envelope
	respMsg := Message{
		Type:    typ,                           // The client will use this type
		Payload: json.RawMessage(payloadBytes), // Pass the marshaled bytes
	}

	// Marshal the final envelope
	msgBytes, err := json.Marshal(respMsg)
	if err != nil {
		log.Printf("Error marshaling envelope: %v", err)
		return
	}

	// Send the message
	c.send <- msgBytes
}
