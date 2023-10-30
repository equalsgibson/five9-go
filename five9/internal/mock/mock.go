package mock

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/equalsgibson/five9-go/five9"
	"github.com/equalsgibson/five9-go/five9/five9types"
)

var GoodCreds = five9types.PasswordCredentials{
	Username: "good-username",
	Password: "good-password",
}

func NewMock(creds five9types.PasswordCredentials) *Mock {
	m := &Mock{
		clientQueue: make(chan []byte),
		serverQueue: make(chan []byte),
	}

	m.Five9 = five9.NewService(creds, five9.SetRoundTripper(m), five9.SetWebsocketHandler(m))

	return m
}

func (mock *Mock) StartMock(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	StartWebsocketError := make(chan error)

	go func() {
		StartWebsocketError <- mock.Five9.Supervisor().StartWebsocket(ctx)
	}()

	// Login complete
	select {
	case <-ctx.Done():
		return errors.New("timed out checking for ping")
	case err := <-StartWebsocketError:
		return err
	case <-mock.loginCompleted:
		break
	}

	// Ping check
	select {
	case <-ctx.Done():
		return errors.New("timed out checking for ping")
	case err := <-StartWebsocketError:
		return err
	case <-mock.pingCompleted:
		break
	}

	// TODO: send message for the data we expect to get
	mock.WriteToClient(ctx, nil)
	mock.WriteToClient(ctx, nil)
	mock.WriteToClient(ctx, nil)
	mock.WriteToClient(ctx, nil)
	mock.WriteToClient(ctx, "messages complete")

	// Confirm last message was processed correctly
	for {
		select {
		case <-ctx.Done():
			return errors.New("timed out checking for ping")
		case err := <-StartWebsocketError:
			return err
		default:
			_, err := mock.Five9.Supervisor().WSAgentState(ctx)
			if errors.Is(err, five9.ErrWebSocketCacheNotReady) {
				continue
			}
			// TODO: verify the incremental came through
			return err
		}
	}

	// Check ping was received
	// setup channel for any error
	return nil
}

type Mock struct {
	Five9              *five9.Service
	loginState         string
	sessionCookieValue string
	username           string
	password           string

	pingCompleted     chan bool
	ConnectionError   error
	clientQueue       chan []byte
	serverQueue       chan []byte
	checkFrameContent func(data []byte)
}

// Roundtrip is the "mock" responses from the server.
func (mock *Mock) RoundTrip(r *http.Request) (*http.Response, error) {
	router := map[string]http.Handler{
		"/login": http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 	// generate cookie value
			// 	// save to internal field
			// 	// set set-cooke header
		}),
	}

	w := httptest.NewRecorder()

	handler, ok := router[r.URL.Path]
	if ok {
		handler.ServeHTTP(w, r)
	} else {
		http.NotFound(w, r)
	}

	return w.Result(), nil
}

func (mock *Mock) Connect(ctx context.Context, connectionURL string, httpClient *http.Client) error {
	if mock.ConnectionError != nil {
		return mock.ConnectionError
	}

	return nil
}

func (mock *Mock) Read(ctx context.Context) ([]byte, error) {
	newMessage := <-mock.clientQueue
	mock.checkFrameContent(newMessage)
	return newMessage, nil
}

func (mock *Mock) Write(ctx context.Context, data []byte) error {
	newMessage := string(data)

	switch newMessage {
	case "ping":
		fileBytes, err := os.ReadFile("test/webSocketFrames/1202_pong.json")
		if err != nil {
			return err
		}

		mock.WriteToClient(ctx, fileBytes)
		mock.pingCompleted <- true
		return nil
	default:
		return errors.New("unsupported message")
	}
}

// Mock the data writes to the WS Server.
func (mock *Mock) WriteToClient(_ context.Context, data []byte) {
	mock.clientQueue <- data
}

func (h *Mock) Close() {}
