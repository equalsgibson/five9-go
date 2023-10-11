package five9

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
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
		response := websocketMessage{
			Context: struct {
				EventID eventID `json:"eventId"`
			}{
				EventID: "1202",
			},
			Payload: "pong",
		}
		message, err := json.Marshal(response)
		h.WriteToClient(ctx, message)
		return err

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
			targetFrame := websocketMessage{}
			if err := json.Unmarshal(data, &targetFrame); err != nil {
				testErr <- err
			}

			if targetFrame.Context.EventID == eventIDPongReceived {
				testErr <- nil
			}
		},
	}

	mockRoundTripper := MockRoundTripper{
		Func: generateWSLoginRequestFuncs(t),
	}

	s := NewService(
		PasswordCredentials{},
		setWebsocketHandler(mockWebsocket),
		setRoundTripper(&mockRoundTripper),
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
	ctx := context.Background()
	testErr := make(chan error)

	mockWebsocket := &MockWebsocketHandler{
		clientQueue: make(chan []byte),
		serverQueue: make(chan []byte),
	}

	mockRoundTripper := MockRoundTripper{
		Func: generateWSLoginRequestFuncs(t),
	}

	s := NewService(
		PasswordCredentials{},
		setWebsocketHandler(mockWebsocket),
		setRoundTripper(&mockRoundTripper),
	)

	go func() {
		if err := s.Supervisor().StartWebsocket(ctx); err != nil {
			testErr <- err

			return
		}
	}()

	mockWebsocket.WriteToClient(ctx, createByteSliceFromFile(t, "test/webSocketFrames/1010_successfulWebSocketConnection.json"))
	mockWebsocket.WriteToClient(ctx, createByteSliceFromFile(t, "test/webSocketFrames/5000_supervisorStats.json"))

	// TODO: maybe need a small sleep here
	time.Sleep(time.Second)

	go func() {
		agents, err := s.Supervisor().AgentState(ctx)
		if err != nil {
			testErr <- err
		}

		if len(agents) != 2 {
			testErr <- fmt.Errorf("expected 2 agents in internal cache, got %d", len(agents))
		}

		// Unblock the channel
		testErr <- nil
	}()

	if err := <-testErr; err != nil {
		t.Fatal(err)
	}
}

func createIoReadCloserFromFile(t *testing.T, filePath string) io.ReadCloser {
	t.Helper()

	file, err := os.Open(filePath)
	if err != nil {
		t.Fatalf("Response Body File Not Found: %s", filePath)
	}

	return io.NopCloser(file)
}

func createByteSliceFromFile(t *testing.T, filePath string) []byte {
	t.Helper()

	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("File Not Found: %s", filePath)
	}

	return fileBytes
}

// The below requests run in order when first starting the websocket service.
func generateWSLoginRequestFuncs(t *testing.T) []func(r *http.Request) (*http.Response, error) {
	t.Helper()

	return []func(r *http.Request) (*http.Response, error){
		func(r *http.Request) (*http.Response, error) { // https://app.five9.com/supsvcs/rs/svc/auth/login
			return &http.Response{
				Body:       createIoReadCloserFromFile(t, "test/supervisorLogin_200.json"),
				StatusCode: http.StatusOK,
			}, nil
		},
		func(r *http.Request) (*http.Response, error) { // supsvcs/rs/svc/supervisors/:userID/login_state
			return &http.Response{
				Body:       createIoReadCloserFromFile(t, "test/loginState_selectStation.json"),
				StatusCode: http.StatusOK,
			}, nil
		},
		func(r *http.Request) (*http.Response, error) { // supsvcs/rs/svc/supervisors/:userID/session_start?force=true
			return &http.Response{
				Body:       http.NoBody,
				StatusCode: http.StatusNoContent,
			}, nil
		},
		func(r *http.Request) (*http.Response, error) { // supsvcs/rs/svc/orgs/:organizationID/users
			return &http.Response{
				Body:       createIoReadCloserFromFile(t, "test/supervisor_getAllUsers_200.json"),
				StatusCode: http.StatusOK,
			}, nil
		},
		func(r *http.Request) (*http.Response, error) { // request_full_statistics
			return &http.Response{
				Body:       http.NoBody,
				StatusCode: http.StatusNoContent,
			}, nil
		},
	}
}
