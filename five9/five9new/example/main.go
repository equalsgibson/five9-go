package main

import (
	"context"
	"log"
	"os"

	"github.com/equalsgibson/five9-go/five9/five9new"
	"github.com/equalsgibson/five9-go/five9/five9types"
)

func main() {
	ctx := context.Background()

	s, err := five9new.NewService(five9types.PasswordCredentials{
		Username: os.Getenv("FIVE9USERNAME"),
		Password: os.Getenv("FIVE9PASSWORD"),
	}, nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	errChan := make(chan error)

	go func() {
		if err := s.StartWebsocket(ctx); err != nil {
			errChan <- err
		}
	}()

	if err := <-errChan; err != nil {
		log.Fatal(err)
	}
}
