package main

import (
	"context"
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
		log.Fatalf("Some error occurred. Err: %s", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

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

	// Start a websocket connection and retry on error
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if err := c.Supervisor().StartWebsocket(ctx); err != nil {
					log.Printf("Websocket exiting, restarting. Here is the error message: %s", err.Error())
					return
				}
			}
		}
	}()

	go func() {
		for range time.NewTicker(time.Second * 2).C {
			agents, err := c.Supervisor().WSAgentState(ctx)
			if err != nil {
				continue
			}

			log.Printf("Found %d agents", len(agents))
		}
	}()
	time.Sleep(time.Second * 10)
	log.Print("Cancelling the context...")
	cancel()
	time.Sleep(time.Second * 2)
}
