package five9_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/equalsgibson/five9-go/five9"
	"github.com/equalsgibson/five9-go/five9/five9types"
)

type MockRoundTripper struct {
	Func []func(r *http.Request) (*http.Response, error)
}

func (mock *MockRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	if len(mock.Func) == 0 {
		return nil, errors.New("end of queue")
	}

	response := mock.Func[0]
	mock.Func = mock.Func[1:]

	return response(r)
}

type MockWebsocketHandler struct {
	ConnectionError   error
	clientQueue       chan []byte
	serverQueue       chan []byte
	checkFrameContent func(data []byte)
}

func (h *MockWebsocketHandler) Connect(ctx context.Context, connectionURL string, httpClient *http.Client) error {
	if h.ConnectionError != nil {
		return h.ConnectionError
	}

	return nil
}

func (h *MockWebsocketHandler) Read(ctx context.Context) ([]byte, error) {
	newMessage := <-h.clientQueue
	h.checkFrameContent(newMessage)
	return newMessage, nil
}

func (h *MockWebsocketHandler) Write(ctx context.Context, data []byte) error {
	newMessage := string(data)

	switch newMessage {
	case "ping":
		fileBytes, err := os.ReadFile("test/webSocketFrames/1202_pong.json")
		if err != nil {
			return err
		}

		h.WriteToClient(ctx, fileBytes)
		return nil

	default:
		return errors.New("unsupported message")
	}
}

// Mock the data writes to the WS Server.
func (h *MockWebsocketHandler) WriteToClient(_ context.Context, data []byte) {
	h.clientQueue <- data
}

func (h *MockWebsocketHandler) Close() {}

func Test_WebSocketPingResponse_Success(t *testing.T) {
	ctx := context.Background()
	testErr := make(chan error)

	mockWebsocket := &MockWebsocketHandler{
		clientQueue: make(chan []byte),
		serverQueue: make(chan []byte),
		checkFrameContent: func(data []byte) {
			targetFrame := five9types.WebsocketMessage{}
			if err := json.Unmarshal(data, &targetFrame); err != nil {
				testErr <- err
			}

			if targetFrame.Context.EventID == five9types.EventIDPongReceived {
				testErr <- nil
			}
		},
	}

	mockRoundTripper := MockRoundTripper{
		Func: generateWSLoginRequestFuncs(t),
	}

	s := five9.NewService(
		five9types.PasswordCredentials{},
		five9.SetWebsocketHandler(mockWebsocket),
		five9.SetRoundTripper(&mockRoundTripper),
	)

	go func() {
		if err := s.Supervisor().StartWebsocket(ctx); err != nil {
			testErr <- err

			return
		}
	}()

	if err := <-testErr; err != nil {
		t.Fatal(err)
	}
}

func Test_GetInternalCache_Success(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	testErr := make(chan error)

	mockWebsocket := &MockWebsocketHandler{
		clientQueue:       make(chan []byte),
		serverQueue:       make(chan []byte),
		checkFrameContent: func(data []byte) {},
	}

	mockRoundTripper := MockRoundTripper{
		Func: generateWSLoginRequestFuncs(t),
	}

	s := five9.NewService(
		five9types.PasswordCredentials{},
		five9.SetWebsocketHandler(mockWebsocket),
		five9.SetRoundTripper(&mockRoundTripper),
		five9.AddRequestPreprocessor(func(r *http.Request) error {
			t.Logf("five9 Rest API Call: [%s] %s", r.Method, r.URL.String())

			return nil
		}),
	)

	go func() {
		if err := s.Supervisor().StartWebsocket(ctx); err != nil {
			testErr <- err
			cancel()
			return
		}
	}()

	mockWebsocket.WriteToClient(ctx, createByteSliceFromFile(t, "test/webSocketFrames/1010_successfulWebSocketConnection.json"))
	t.Logf("Got to here")
	mockWebsocket.WriteToClient(ctx, createByteSliceFromFile(t, "test/webSocketFrames/5000_stats.json"))
	t.Logf("Got to here 2")
	go func() {
		maxAttempts := 2
		for i := 0; i <= maxAttempts; i++ {
			agents, err := s.Supervisor().AgentState(ctx)
			if err != nil {
				testErr <- err
				cancel()
			}

			if len(agents) != 2 {
				if i != maxAttempts {
					time.Sleep(time.Millisecond)
					// Try again
					continue
				}
				testErr <- fmt.Errorf("expected 2 agents in internal cache, got %d", len(agents))
				cancel()
			}
			// Unblock the channel
			testErr <- nil
			cancel()
		}
	}()

	for {
		select {
		case err := <-testErr:
			if err != nil {
				t.Fatal(err)
			}
		case <-ctx.Done():
			return
		}
	}
}
