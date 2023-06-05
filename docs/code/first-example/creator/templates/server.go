package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
)

var PORT string = ":8080"

func main() {
	// create echo server
	e := echo.New()
	e.HideBanner = true

	// add middleware
	err := setupMiddleware(e)
	if err != nil {
		panic(err)
	}

	// setup router
	err = setupRouter(e)
	if err != nil {
		panic(err)
	}

	//
	// code to run server and enable graceful shutdown
	//
	// Start server with background goroutine
	go func() {
		if err := e.Start(PORT); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	// wait on a quit signal
	<-quit

	// start the shutdown process
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
