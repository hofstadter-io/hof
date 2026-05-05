package runtime

import (
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/gorilla/websocket"
	"github.com/hofstadter-io/hof/lib/consts"
	"github.com/labstack/echo/v4"
)

func (r *Runtime) Run() error {
	port := ":2257"

	routes := r.e.Routes()
	sort.Slice(routes, func(i, j int) bool {
		lhs, rhs := routes[i], routes[j]
		if lhs.Path < rhs.Path {
			return true
		}
		if lhs.Path > rhs.Path {
			return false
		}
		return lhs.Method < rhs.Method
	})
	for _, route := range routes {
		fmt.Printf("%-6s %s\n", route.Method, route.Path)
	}

	go r.runRegistrar()
	r.e.Logger.Fatal(r.e.Start(port))

	return nil
}

func (r *Runtime) runRegistrar() {
	for {
		select {
		case client := <-r.register:
			r.mu.Lock()
			r.clients[client] = true
			log.Printf("Client connected. Total clients: %d", len(r.clients))
			r.mu.Unlock()

		case client := <-r.unregister:
			r.mu.Lock()
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
				close(client.send)
				log.Printf("Client disconnected. Total clients: %d", len(r.clients))
			}
			r.mu.Unlock()

		}
	}
}

var upgrader = websocket.Upgrader{
	// Allow all origins for local development.
	// In production, you might want to restrict this.
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (R *Runtime) serveWs(c echo.Context) error {
	// TODO, we need to pull the user info / auth from here before upgrading and such
	// TODO, store user info on the client type
	user := c.Request().Header.Get(consts.VEG_USER_HEADER)
	if user == "" {
		user = consts.VEG_DEFAULT_USER
	}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return fmt.Errorf("while upgrading websocket in runtime")
	}

	client := &Client{
		User:          user,
		conn:          conn,
		send:          make(chan []byte, 256), // 256-message buffer
		handleMessage: R.handleMessage,
		State:         make(map[string]any),
	}

	// Register the new client with the hub
	R.register <- client

	// Start the read/write pumps as goroutines
	go client.writePump()
	go R.readPump(client)

	return nil
}
