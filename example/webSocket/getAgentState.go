package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/equalsgibson/five9-go/five9"
	"github.com/equalsgibson/five9-go/five9/five9types"
)

func startConnection(ctx context.Context, c *five9.Service) error {
	if err := c.Supervisor().StartWebsocket(ctx); err != nil {
		if !errors.Is(err, context.Canceled) {
			return fmt.Errorf("websocket exiting, restarting. Here is the error message: %s", err.Error())
		}
	}

	return nil
}

func main() {
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

	// Run a function every 5 seconds to obtain some information from the
	// supervisor websocket connection.

	go func() {
		for range time.NewTicker(time.Second * 5).C {
			agents, err := c.Supervisor().WSAgentState(ctx)
			if err != nil {
				log.Printf("Err: %s", err)

				continue
			}

			log.Println("+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+")
			log.Printf("Found %d total agents", len(agents))
			log.Println("+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+")
		}
	}()

	count := runtime.NumGoroutine()
	defer func() {
		time.Sleep(time.Second)
		diff := runtime.NumGoroutine() - count
		pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
		if diff > 0 {
			log.Printf("Too many goruoutes: %d extra (%d now and started with %d)", diff, runtime.NumGoroutine(), count)
		} else {
			log.Printf("Looking good")
		}
	}()

	errCount := 0
	for {
		if err := startConnection(ctx, c); err != nil {
			errCount++
			log.Print(err)
			time.Sleep(time.Second * 1)
		}

		if errCount > 3 {
			return
		}
	}
}
