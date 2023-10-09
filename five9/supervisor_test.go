package five9

import (
	"context"
	"encoding/json"
	"errors"
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
	ConnectionError error
	clientQueue     chan []byte
	serverQueue     chan []byte
}

func (h *MockWebsocketHandler) Connect(ctx context.Context, connectionURL string, httpClient *http.Client) error {
	if h.ConnectionError != nil {
		return h.ConnectionError
	}

	return nil
}

func (h *MockWebsocketHandler) Read(ctx context.Context) ([]byte, error) {
	newMessage := <-h.clientQueue

	return newMessage, nil
}

func (h *MockWebsocketHandler) Write(_ context.Context, data []byte) error {
	h.serverQueue <- data

	return nil
}

// Mock the data writes to the WS Server
func (h *MockWebsocketHandler) WriteToClient(_ context.Context, data []byte) {
	h.clientQueue <- data
}

// Mock the responses from the WS Server
func (h *MockWebsocketHandler) ReadFromServer(ctx context.Context) ([]byte, error) {
	newMessageBytes := <-h.serverQueue
	newMessageString := string(newMessageBytes)

	switch newMessageString {
	case "ping":
		response := websocketMessage{
			Context: struct {
				EventID eventID "json:\"eventId\""
			}{
				EventID: "1202",
			},
			Payload: "pong",
		}
		return json.Marshal(response)

	default:
		return nil, errors.New("unsupported message")
	}
}

func (h *MockWebsocketHandler) Close() {}

func Test_WebSocketPingResponse_Success(t *testing.T) {
	ctx := context.Background()

	mockWebsocket := &MockWebsocketHandler{
		clientQueue: make(chan []byte),
		serverQueue: make(chan []byte),
	}

	{ // Spin the mock WS Server up, and process messages
		go func() {
			for {
				messageBytes, err := mockWebsocket.ReadFromServer(ctx)
				if err != nil {
					//TODO: Panics are not OK - they do not fail the test.
					panic(err)
				}
				mockWebsocket.WriteToClient(ctx, messageBytes)
			}
		}()
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
			panic(err)
		}
	}()

	// Sleep for 6 seconds to allow a ping to be sent
	// TODO: change package to send a ping on connect
	time.Sleep(time.Second * 6)
}

func Test_GetInternalCache_Success(t *testing.T) {
	ctx := context.Background()

	mockWebsocket := &MockWebsocketHandler{
		clientQueue: make(chan []byte),
		serverQueue: make(chan []byte),
	}

	{ // Spin the mock WS Server up, and process messages
		go func() {
			for {
				messageBytes, err := mockWebsocket.ReadFromServer(ctx)
				if err != nil {
					panic(err)
				}
				mockWebsocket.WriteToClient(ctx, messageBytes)
			}
		}()
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
			panic(err)
		}
	}()

	mockWebsocket.WriteToClient(ctx, createByteSliceFromFile(t, "test/webSocketFrames/1010_successfulWebSocketConnection.json"))
	mockWebsocket.WriteToClient(ctx, createByteSliceFromFile(t, "test/webSocketFrames/5000_supervisorStats.json"))

	// TODO: maybe need a small sleep here
	time.Sleep(time.Second)
	agents, err := s.Supervisor().AgentState(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(agents) != 2 {
		t.Fatalf("expected 2 agents in internal cache")
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
