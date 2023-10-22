package main

import (
	"context"
	"log"
	"net/http"
	"os"

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

	domainUsers, err := c.Supervisor().GetAllDomainUsers(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Print a count of the users within your Five9 Domain
	log.Printf("You have %d users within your Five9 Domain.\n", len(domainUsers))
}
