package internal_test

import (
	"context"
	"log"
	"net/http"
	"testing"

	"github.com/equalsgibson/five9-go/five9"
	"github.com/equalsgibson/five9-go/five9/five9types"
	"github.com/equalsgibson/five9-go/five9/internal"
)

func TestC(t *testing.T) {
	ctx := context.Background()

	testServer := internal.NewMockFive9Server(t)
	defer testServer.Close()

	log.Println(testServer.URL)

	service := five9.NewService(five9types.PasswordCredentials{}, five9.SetTestServerLoginURL(testServer.URL), five9.AddRequestPreprocessor(func(r *http.Request) error {
		log.Printf("five9 Rest API Call: [%s] %s", r.Method, r.URL.String())

		return nil
	}))
	e, err := service.Supervisor().GetOwnUserInfo(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(e)
	log.Println("Got here")
}
