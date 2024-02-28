package five9new_test

import (
	"context"
	"errors"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/equalsgibson/five9-go/five9/five9new"
	"github.com/equalsgibson/five9-go/five9/five9types"
)

type testRoundTripper struct {
	handler func(r *http.Request) (*http.Response, error)
}

func (t testRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	return t.handler(r)
}

func TestService(t *testing.T) {
	_, server := net.Pipe()
	ctx := context.Background()
	testService := createTestService(
		t,
		ctx,
		testRoundTripper{
			handler: func(r *http.Request) (*http.Response, error) {

				return nil, nil
			},
		},
		server,
	)

	testService.StartWebsocket(ctx)
}

func TestService_SessionExpiring(t *testing.T) {
}

func TestService_ServiceMigrated(t *testing.T) {

}
func TestService_PendingMaintenanceNotice(t *testing.T) {
}

func TestService_PongFailure(t *testing.T) {
}

func TestService_HandleRetryableError(t *testing.T) {

}

func TestService_HandleFatalError(t *testing.T) {

}

func TestService_BadCredentials(t *testing.T) {
	// testService := createTestService(t, http.DefaultTransport, nil)

	// testService.StartWebsocket()
}

func createTestService(
	t *testing.T,
	parentContext context.Context,
	roundTripper http.RoundTripper,
	testServerConn net.Conn,
) *five9new.Service {
	ctx, cancel := context.WithCancelCause(parentContext)
	socketConnected := make(chan bool)
	service, err := five9new.NewService(five9types.PasswordCredentials{
		Username: "fsd",
		Password: "s",
	}, testRoundTripper{
		handler: func(r *http.Request) (*http.Response, error) {
			if r.URL.Path == "/connection/to/websocket" {
				// testServerConn.RemoteAddr()
				// return websocket conn details
				socketConnected <- true
				return nil, nil
			}

			// Login Endpoint
			// check password

			return roundTripper.RoundTrip(r)
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		timeout := time.NewTicker(time.Second * 5)

		select {
		case <-socketConnected:
			time.Sleep(time.Second)
			// send off time based frame queue
		case <-timeout.C:
			cancel(errors.New("didn't connect withine timeframe"))
		case <-ctx.Done():
			return
		}
	}()

	return service
}

/*
	Cache:
		LastControlReceived:time.now.add(-61)

		KeepaliveTimeout = 60
		Conn
*/

// http < - > conn

// conn <-

/*
2 types of frames
- response frames (trigger based)
--- I just sent a ping, send a pong in response
- timebased frames (this could be a QUEUE of frames that get sent in order)
--- Scheduled frames -> full update, then incremental add, then incremental remove, then add etc.
*/

// // TestWebsocketServer {
// ValidPassweord: Token
// }
