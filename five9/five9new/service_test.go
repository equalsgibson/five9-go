package five9new_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/equalsgibson/five9-go/five9/five9new"
	"github.com/equalsgibson/five9-go/five9/five9new/internal/study"
	"github.com/equalsgibson/five9-go/five9/five9types"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func TestService_Success(t *testing.T) {
	_, server := net.Pipe()

	ctx := context.Background()

	mockFive9Server := &TestFive9HTTPServer{
		serverConn: &server,
		t:          t,
	}

	_ = createTestService(
		t,
		ctx,
		mockFive9Server,
		server,
	)
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
}

func createTestService(
	t *testing.T,
	parentContext context.Context,
	testFive9HTTPServer *TestFive9HTTPServer,
	testServerConn net.Conn,
) *five9new.Service {
	ctx, cancel := context.WithCancelCause(parentContext)
	socketConnected := make(chan bool)

	service, err := five9new.NewService(
		five9types.PasswordCredentials{
			Username: "fsd",
			Password: "s",
		},
		testFive9HTTPServer,
		func(ctx context.Context, network, addr string) (net.Conn, error) {
			socketConnected <- true
			return testServerConn, nil
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	go func() {
		if err := service.StartWebsocket(ctx); err != nil {
			cancel(err)
		}
	}()

	go func() {
		<-socketConnected

		if err := testFive9HTTPServer.writeQueue(); err != nil {
			cancel(err)
		}
	}()

	timeout := time.NewTicker(time.Second * 5)

	select {
	case <-socketConnected:
		time.Sleep(time.Second)
		// send off time based frame queue
	case <-timeout.C:
		cancel(errors.New("didn't connect withine timeframe"))
	case <-ctx.Done():
		t.Fatal(context.Cause(ctx))
	}

	return service
}

type TestFive9HTTPServer struct {
	validPassword string
	userState     UserState
	serverConn    *net.Conn
	t             *testing.T
	PingCount     int
}

func (t *TestFive9HTTPServer) RoundTrip(r *http.Request) (*http.Response, error) {
	log.Println(r.URL.Path)
	if r.URL.Path == "/supsvcs/rs/svc/auth/login" {
		return t.loginHandler(r)
	}
	return nil, nil
}

func (t *TestFive9HTTPServer) loginHandler(r *http.Request) (*http.Response, error) {
	return study.ServeAndValidate(
		t.t,
		&study.TestResponseFile{
			StatusCode:        200,
			FilePath:          "test/auth.json",
			ResponseModifiers: nil,
		},
		study.ExpectedTestRequest{
			Method: http.MethodPost,
			Path:   "/supsvcs/rs/svc/auth/login",
		},
	)(t.t, r)
}

func (t *TestFive9HTTPServer) readWSMessage() error {
	log.Println("Reading a SERVER frame")
	for {
		header, err := ws.ReadHeader(*t.serverConn)
		if err != nil {
			return err
			// handle error
		}

		if header.OpCode == ws.OpPing {
			t.PingCount++
			if t.PingCount > 2 {
				continue
			}
			t.write(nil, ws.OpPong)
			continue
		}

		payload := make([]byte, header.Length)
		_, err = io.ReadFull(*t.serverConn, payload)
		if err != nil {
			return err
			// handle error
		}
		if header.Masked {
			ws.Cipher(payload, header.Mask, 0)
		}

		log.Println("SERVER Payload; ", string(payload))

		continue
	}
}

func (s *TestFive9HTTPServer) write(
	payload any,
	opCode ws.OpCode,
) error {
	writer := wsutil.NewWriter(*s.serverConn, ws.StateServerSide, opCode)
	encoder := json.NewEncoder(writer)

	if err := encoder.Encode(payload); err != nil {
		return err
	}

	return writer.Flush()
}

type UserState struct {
	UserID    string
	SessionID string
}

type TestFrame struct {
	Payload any
	OpCode  ws.OpCode
	Delay   time.Duration
}

func (t *TestFive9HTTPServer) writeQueue() error {
	queue := []TestFrame{
		TestFrame{
			Payload: "Hello",
			OpCode:  ws.OpText,
			Delay:   time.Second,
		},
	}

	for _, q := range queue {
		time.Sleep(q.Delay)
		if err := t.write(q.Payload, q.OpCode); err != nil {
			return err
		}
	}

	return nil
}
