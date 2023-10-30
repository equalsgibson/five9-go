package five9_test

import (
	"context"
	"testing"

	"github.com/equalsgibson/five9-go/five9/internal/mock"
)

func Test_WebSocketPingResponse_Success(t *testing.T) {
	ctx := context.Background()

	m := mock.NewMock(mock.GoodCreds)
	// configure the mock for this test, no config need for success path

	// Check to make sure setup is complete (all incremental updates received)
	if err := m.StartMock(ctx); err != nil {
		t.Fatal(err)
	}

	agents, err := m.Five9.Supervisor().WSAgentState(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// TODO assert agent list
}

// func Test_GetInternalCache_Success(t *testing.T) {
// 	ctx := context.Background()
// 	testErr := make(chan error)

// 	mockWebsocket := &MockWebsocketHandler{
// 		clientQueue:       make(chan []byte),
// 		serverQueue:       make(chan []byte),
// 		checkFrameContent: func(data []byte) {},
// 	}

// 	mockRoundTripper := MockRoundTripper{
// 		Func: generateWSLoginRequestFuncs(t),
// 	}

// 	s := five9.NewService(
// 		five9types.PasswordCredentials{},
// 		five9.SetWebsocketHandler(mockWebsocket),
// 		five9.SetRoundTripper(&mockRoundTripper),
// 		five9.AddRequestPreprocessor(func(r *http.Request) error {
// 			t.Logf("API Call Made: [%s] %s\n", r.Method, r.URL.String())

// 			return nil
// 		}),
// 	)

// 	go func() {
// 		if err := s.Supervisor().StartWebsocket(ctx); err != nil {
// 			testErr <- err

// 			return
// 		}
// 	}()

// 	mockWebsocket.WriteToClient(ctx, createByteSliceFromFile(t, "test/webSocketFrames/1010_successfulWebSocketConnection.json"))
// 	mockWebsocket.WriteToClient(ctx, createByteSliceFromFile(t, "test/webSocketFrames/5000_stats.json"))
// 	time.Sleep(time.Second)

// 	agents, err := s.Supervisor().WSAgentState(ctx)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if len(agents) != 2 {
// 		t.Fatalf("expected 2 agents in internal cache, got %d", len(agents))
// 	}

// 	if err := <-testErr; err != nil {
// 		t.Fatal(err)
// 	}
// }
