package five9

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"sync"
	"time"
)

type client struct {
	httpClient           *http.Client
	credentials          PasswordCredentials
	requestPreProcessors []func(r *http.Request) error
	login                *loginResponse
	loginMutex           *sync.Mutex
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

func (c *client) requestWithAuthentication(request *http.Request, target any) error {
	login, err := c.getLogin(request.Context())
	if err != nil {
		return err
	}

	request.URL.Scheme = "https"
	request.URL.Host = login.GetAPIHost()
	request.URL.Path = strings.ReplaceAll(request.URL.Path, ":userID", string(login.UserID))
	request.URL.Path = strings.ReplaceAll(request.URL.Path, ":organizationID", string(login.OrgID))

	return c.request(request, target)
}

func (c *client) getLogin(ctx context.Context) (*loginResponse, error) {
	{ // check for existing login
		if c.login != nil {
			return c.login, nil
		}

		c.loginMutex.Lock()
		defer c.loginMutex.Unlock()

		if c.login != nil {
			return c.login, nil
		}
	}

	login, err := c.endpointLogin(ctx)
	if err != nil {
		return nil, err
	}

	c.login = &login

	if err := c.endpointGetSessionMetadata(ctx); err != nil {
		return nil, err
	}

	loginState, err := c.endpointGetLoginState(ctx)
	if err != nil {
		return nil, err
	}

	if loginState == "SELECT_STATION" {
		if err := c.endpointStartSession(ctx); err != nil {
			return nil, err
		}
	}

	return c.login, nil
}

func (c *client) endpointLogin(ctx context.Context) (loginResponse, error) {
	// Create new, blank cookiejar to make sure old cookies are not used
	newJar, _ := cookiejar.New(nil)
	c.httpClient.Jar = newJar

	payload := loginPayload{
		PasswordCredentials: c.credentials,
		AppKey:              "web-ui",
		Policy:              PolicyAttachExisting,
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://app.five9.com/supsvcs/rs/svc/auth/login",
		structToReaderCloser(payload),
	)
	if err != nil {
		return loginResponse{}, err
	}

	target := loginResponse{}

	if err := c.request(request, &target); err != nil {
		return loginResponse{}, err
	}

	return target, nil
}

func (c *client) endpointGetSessionMetadata(ctx context.Context) error {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"/supsvcs/rs/svc/auth/metadata",
		http.NoBody,
	)
	if err != nil {
		return err
	}

	return c.requestWithAuthentication(request, nil)
}

func (c *client) endpointGetLoginState(ctx context.Context) (userLoginState, error) {
	var target userLoginState

	tries := 0
	for tries < 3 {
		tries++

		request, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			"/supsvcs/rs/svc/supervisors/:userID/login_state",
			http.NoBody,
		)
		if err != nil {
			return "", err
		}

		if err := c.requestWithAuthentication(request, &target); err != nil {
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

func (c *client) endpointStartSession(ctx context.Context) error {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		"/supsvcs/rs/svc/supervisors/:userID/session_start",
		structToReaderCloser(StationInfo{
			StationID:   "",
			StationType: "EMPTY",
			Force:       true,
		}),
	)
	if err != nil {
		return err
	}

	if err := c.requestWithAuthentication(request, nil); err != nil {
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
