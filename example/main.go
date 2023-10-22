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
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Fatalf("Could not get environment variables. Err: %s", err)
	}

	// ctx, cancel := context.WithCancel(context.Background())
	ctx := context.Background()

	c := five9.NewService(
		five9types.PasswordCredentials{
			Username: os.Getenv("USERNAME"),
			Password: os.Getenv("PASSWORD"),
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

	// TODO: Make a comment explaining logic
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			agents, err := c.Supervisor().WSAgentState(ctx)
			if err != nil {
				continue
			}
			log.Printf("Found %d agents", len(agents))
		}
	}
}
