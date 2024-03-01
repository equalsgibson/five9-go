package five9new

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/equalsgibson/five9-go/five9/five9types"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type Service struct {
	client           *http.Client
	credentials      *five9types.PasswordCredentials
	keepaliveTimeout time.Duration
	netDialer        func(ctx context.Context, network, addr string) (net.Conn, error)
}

func (s *Service) StartWebsocket(ctx context.Context) error {
	// Login will set the cookies in the cookiejar
	loginResponse, err := s.login(ctx, *s.credentials)
	if err != nil {
		return err
	}

	log.Println(loginResponse)

	connectionURL, err := url.Parse("https://app.five9.com")
	if err != nil {
		return err
	}

	authorizationCookies := s.client.Jar.Cookies(connectionURL)

	headers := ws.HandshakeHeaderHTTP{}
	for _, cookie := range authorizationCookies {
		headers["Cookie"] = append(headers["Cookie"], cookie.String())
	}

	log.Printf("%+v\n", headers)

	dialer := ws.Dialer{
		Header:  headers,
		NetDial: s.netDialer,
	}

	conn, _, _, err := dialer.Dial(ctx, "wss://app-atl.five9.com/supsvcs/sws/ws")
	if err != nil {
		return err
	}

	errorChan := make(chan error)
	go func() {
		if err := s.ping(ctx, conn); err != nil {
			errorChan <- err
		}
	}()

	go func() {
		if err := s.read(conn); err != nil {
			errorChan <- err
		}
	}()

	return <-errorChan
}

func (s *Service) GetACDStatus() ([]five9types.ACDState, error) {
	return nil, nil
}

func (s *Service) GetAgentStatus() ([]five9types.AgentState, error) {
	return nil, nil
}

func NewService(
	credentials five9types.PasswordCredentials,
	roundtripper http.RoundTripper,
	netDialer func(
		ctx context.Context,
		network string,
		addr string,
	) (net.Conn, error),
) (*Service, error) {
	cookieJar, _ := cookiejar.New(nil)

	client := &http.Client{
		Transport: roundtripper,
		Jar:       cookieJar,
	}

	return &Service{
		client:           client,
		credentials:      &credentials,
		keepaliveTimeout: time.Second * 60,
		netDialer:        netDialer,
	}, nil
}

/* ******************************************
	This is where I tried to implement stuff -
	Warning: Here be dragons
****************************************** */

func (s *Service) ping(ctx context.Context, conn net.Conn) error {
	log.Println("Writing a ping")
	if err := s.write(conn, nil, ws.OpPing); err != nil {
		return err
	}

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case currentPingSentTime := <-ticker.C:
			log.Println("Writing a ping: ", currentPingSentTime)
			if err := s.write(conn, nil, ws.OpPing); err != nil {
				return err
			}

		case <-ctx.Done():
			return nil
		}
	}

}

func (s *Service) read(conn net.Conn) error {
	log.Println("Reading a CLIENT frame")
	for {
		header, err := ws.ReadHeader(conn)
		if err != nil {
			return err
			// handle error
		}

		payload := make([]byte, header.Length)
		_, err = io.ReadFull(conn, payload)
		if err != nil {
			return err
			// handle error
		}
		if header.Masked {
			ws.Cipher(payload, header.Mask, 0)
		}

		log.Println("CLIENT Payload; ", string(payload))

		continue

	}
}

func (s *Service) write(
	conn net.Conn,
	payload any,
	opCode ws.OpCode,
) error {
	log.Println("Writing a CLIENT frame")
	writer := wsutil.NewWriter(conn, ws.StateClientSide, opCode)
	encoder := json.NewEncoder(writer)

	if err := encoder.Encode(payload); err != nil {
		return err
	}

	return writer.Flush()
}

/* ******************************************
	REST for login and REST Helpers
****************************************** */

func (s *Service) login(ctx context.Context, credentials five9types.PasswordCredentials) (five9types.LoginResponse, error) {
	payload := five9types.LoginPayload{
		PasswordCredentials: credentials,
		AppKey:              "web-ui",
		Policy:              five9types.PolicyForceIn,
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://app.five9.com/supsvcs/rs/svc/auth/login",
		structToReaderCloser(payload),
	)
	if err != nil {
		return five9types.LoginResponse{}, err
	}

	target := five9types.LoginResponse{}

	if err := s.request(request, &target); err != nil {
		return five9types.LoginResponse{}, err
	}

	return target, nil

}

func structToReaderCloser(v any) io.Reader {
	vBytes, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	return bytes.NewReader(vBytes)
}

func (s *Service) request(request *http.Request, target any) error {
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")

	response, err := s.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("error, %d", response.StatusCode)
	}

	if target != nil {
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(bodyBytes, target); err != nil {
			return err
		}
	}

	return nil
}
