package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/equalsgibson/five9-go/five9"
	"github.com/equalsgibson/five9-go/five9/five9types"
)

func main() {
	// Set up a new Five9 service
	ctx := context.Background()
	c := five9.NewService(
		five9types.PasswordCredentials{
			Username: os.Getenv("FIVE9USERNAME"),
			Password: os.Getenv("FIVE9PASSWORD"),
		},
		five9.AddRequestPreprocessor(func(r *http.Request) error {
			log.Printf("five9 Rest API Call: [%s] %s", r.Method, r.URL.String())

			return nil
		}),
	)

	// Start a websocket connection
	go func() {
		if err := c.Supervisor().StartWebsocket(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				log.Printf("Websocket exiting, restarting. Here is the error message: %s", err.Error())
			}
		}
	}()

	// Run a function every 5 seconds to obtain a count of Agents
	// from the supervisor websocket connection.
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			agents, err := c.Supervisor().WSAgentState(ctx)
			if err != nil {
				log.Printf("Err: %s", err)
				continue
			}
			log.Printf("Found %d agents", len(agents))
		}
	}
}
