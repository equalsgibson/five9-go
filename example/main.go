package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/equalsgibson/five9-go/five9"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Fatalf("Some error occurred. Err: %s", err)
	}

	ctx := context.Background()

	c := five9.NewService(
		five9.PasswordCredentials{
			Username: os.Getenv("USERNAME"),
			Password: os.Getenv("PASSWORD"),
		},
		five9.AddRequestPreprocessor(func(r *http.Request) error {
			log.Printf("five9 Rest API Call: [%s] %s", r.Method, r.URL.String())

			return nil
		}),
	)

	go func() {
		for {
			if err := c.Supervisor().StartWebsocket(ctx); err != nil {
				log.Printf("Websocket exiting, restarting. Here is the error message: %s", err.Error())
			}
		}
	}()

	for range time.NewTicker(time.Second * 2).C {
		agents, err := c.Supervisor().AgentState(ctx)
		if err != nil {
			continue
		}

		log.Printf("Found %d agents", len(agents))
	}
}
