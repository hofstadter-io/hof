package runtime

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (r *Runtime) Run() error {
	port := ":2257"

	r.setupServer()
	go r.runRegistrar()
	r.e.Logger.Fatal(r.e.Start(port))

	return nil
}

func (r *Runtime) setupServer() error {
	e := echo.New()
	e.HideBanner = true

	// middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// routes
	e.GET("/", r.serveWs)

	e.GET("/alive", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	// TODO metrics & otel

	// save & return
	r.e = e
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
	// TODO, look for session ID, or do we add that to the message?
	// we will likely have multiple sessions on one websocket
	// we just need to pull the user info / auth from here before upgrading and such

	// TODO, store user info on the client type

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return fmt.Errorf("while upgrading websocket in runtime")
	}

	client := &Client{
		User:          "tony",
		conn:          conn,
		send:          make(chan []byte, 256), // 256-message buffer
		handleMessage: R.handleMessage,
	}

	// Register the new client with the hub
	R.register <- client

	// Start the read/write pumps as goroutines
	go client.writePump()
	go R.readPump(client)

	return nil
}
