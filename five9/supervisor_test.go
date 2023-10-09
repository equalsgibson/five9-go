package five9

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"testing"
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
	return nil
}

func (h *MockWebsocketHandler) WriteToClient(_ context.Context, data []byte) {
	h.clientQueue <- data
}

func (h *MockWebsocketHandler) Close() {}

func TestPing(t *testing.T) {
	ctx := context.Background()

	mockWebsocket := &MockWebsocketHandler{
		clientQueue: make(chan []byte),
	}
	mockRoundTripper := MockRoundTripper{
		Func: []func(r *http.Request) (*http.Response, error){
			func(r *http.Request) (*http.Response, error) {
				return &http.Response{
					Body:       createIoReadCloserFromFile(t, "test/supervisorLogin_200.json"),
					StatusCode: http.StatusOK,
				}, nil
			},
			func(r *http.Request) (*http.Response, error) {
				return &http.Response{
					Body:       createIoReadCloserFromFile(t, "test/loginState_selectStation.json"),
					StatusCode: http.StatusOK,
				}, nil
			},
			func(r *http.Request) (*http.Response, error) {
				return &http.Response{
					Body:       http.NoBody,
					StatusCode: http.StatusNoContent,
				}, nil
			},
			func(r *http.Request) (*http.Response, error) {
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
		},
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

	agents, err := s.Supervisor().AgentState(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(agents) == 2 {
		return
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
