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

	apiUserInfo, err := c.Supervisor().GetOwnUserInfo(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Start a websocket connection
	go func() {
		if err := c.Supervisor().StartWebsocket(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				log.Fatalf("Websocket exiting, restarting. Here is the error message: %s", err.Error())
			}
		}
	}()

	// Run a function every 5 seconds to obtain some information from the
	// supervisor websocket connection.
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for range ticker.C {
		agents, err := c.Supervisor().WSAgentState(ctx)
		if err != nil {
			log.Printf("Err: %s", err)
			continue
		}

		myUserState, ok := agents[apiUserInfo.UserName]
		if !ok {
			log.Fatal("could not find API User ID in agent state map")
		}

		log.Println("+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+")
		log.Printf("Found %d total agents", len(agents))
		log.Printf("Found API User ID and UserName: %s -> %s", apiUserInfo.ID, apiUserInfo.UserName)
		log.Printf("Found the API Users Current State: %+v", myUserState.Presence.CurrentState)
		log.Println("+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+")
	}

}
