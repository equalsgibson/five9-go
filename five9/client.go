package five9

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

type client struct {
	httpClient           *http.Client
	credentials          PasswordCredentials
	requestPreProcessors []func(r *http.Request) error
}

const (
	supervisorAPIContextPath = "supsvcs/rs/svc"
	agentAPIContextPath      = "appsvcs/rs/svc"
)

type authenticationState struct {
	client         *client
	loginResponse  *loginResponse
	loginMutex     *sync.Mutex
	apiContextPath string
}

func (a *authenticationState) endpointGetSessionMetadata(ctx context.Context) error {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("/%s/auth/metadata", a.apiContextPath),
		http.NoBody,
	)
	if err != nil {
		return err
	}

	return a.requestWithAuthentication(request, nil)
}

func (a *authenticationState) requestWithAuthentication(request *http.Request, target any) error {
	login, err := a.getLogin(request.Context())
	if err != nil {
		return err
	}

	request.URL.Scheme = "https"
	request.URL.Host = login.GetAPIHost()
	request.URL.Path = strings.ReplaceAll(request.URL.Path, ":userID", string(login.UserID))
	request.URL.Path = strings.ReplaceAll(request.URL.Path, ":organizationID", string(login.OrgID))

	return a.client.request(request, target)
}

func (a *authenticationState) getLogin(
	ctx context.Context,
) (*loginResponse, error) {
	{ // check for existing login
		if a.loginResponse != nil {
			return a.loginResponse, nil
		}

		a.loginMutex.Lock()
		defer a.loginMutex.Unlock()

		if a.loginResponse != nil {
			return a.loginResponse, nil
		}
	}

	login, err := a.endpointLogin(ctx)
	if err != nil {
		return nil, err
	}

	a.loginResponse = &login

	if err := a.endpointGetSessionMetadata(ctx); err != nil {
		return nil, err
	}

	loginState, err := a.endpointGetLoginState(ctx)
	if err != nil {
		return nil, err
	}

	if loginState == "SELECT_STATION" {
		if err := a.endpointStartSession(ctx); err != nil {
			return nil, err
		}
	}

	return a.loginResponse, nil
}

func (a *authenticationState) endpointLogin(ctx context.Context) (loginResponse, error) {
	payload := loginPayload{
		PasswordCredentials: a.client.credentials,
		AppKey:              "web-ui",
		Policy:              PolicyAttachExisting,
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("https://app.five9.com/%s/auth/login", a.apiContextPath),
		structToReaderCloser(payload),
	)
	if err != nil {
		return loginResponse{}, err
	}

	target := loginResponse{}

	if err := a.client.request(request, &target); err != nil {
		return loginResponse{}, err
	}

	return target, nil
}

func (a *authenticationState) endpointGetLoginState(ctx context.Context) (userLoginState, error) {
	path := "agents"
	if a.apiContextPath == supervisorAPIContextPath {
		path = "supervisors"
	}

	var target userLoginState

	tries := 0
	for tries < 3 {
		tries++

		request, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			fmt.Sprintf(
				"/%s/%s/:userID/login_state",
				a.apiContextPath,
				path,
			),
			http.NoBody,
		)
		if err != nil {
			return "", err
		}

		if err := a.requestWithAuthentication(request, &target); err != nil {
			five9Error, ok := err.(*Error)
			if ok && five9Error.StatusCode == http.StatusUnauthorized {
				// The login is not registered by other endpoints for a short time.
				// I think this has to do with Five9 propagating the session across their data centers.
				// We login using the app.five9.com domain but then make subsequent calls to the data center specific domain
				time.Sleep(time.Second * 2)

				continue
			}

			return "", err
		}

		return target, nil
	}

	return "", errors.New("Five9 login timeout")
}

func (c *client) request(request *http.Request, target any) error {
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")

	for _, requestPreProcessor := range c.requestPreProcessors {
		if err := requestPreProcessor(request); err != nil {
			return err
		}
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode >= http.StatusBadRequest {
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		responseErr := &Error{
			StatusCode: response.StatusCode,
			Body:       bodyBytes,
		}

		if err := json.Unmarshal(bodyBytes, responseErr); err != nil {
			return err
		}

		return responseErr
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

func (a *authenticationState) endpointStartSession(ctx context.Context) error {
	path := "agents"
	if a.apiContextPath == supervisorAPIContextPath {
		path = "supervisors"
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		fmt.Sprintf(
			"/%s/%s/:userID/session_start?force=true",
			a.apiContextPath,
			path,
		),
		structToReaderCloser(StationInfo{
			StationID:   "",
			StationType: "EMPTY",
		}),
	)
	if err != nil {
		return err
	}

	if err := a.requestWithAuthentication(request, nil); err != nil {
		return err
	}

	return nil
}

func structToReaderCloser(v any) io.Reader {
	vBytes, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	return bytes.NewReader(vBytes)
}
