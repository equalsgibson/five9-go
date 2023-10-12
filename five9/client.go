package five9

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	agentLogin           *loginResponse
	agentLoginMutex      *sync.Mutex
	supervisorLogin      *loginResponse
	supervisorLoginMutex *sync.Mutex
}

const (
	supervisorAPIContextPath = "supsvcs/rs/svc"
	agentAPIContextPath      = "appsvcs/rs/svc"
)

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

func (c *client) requestWithSupervisorAuthentication(request *http.Request, target any) error {
	login, err := c.getSupervisorLogin(request.Context())
	if err != nil {
		return err
	}

	request.URL.Scheme = "https"
	request.URL.Host = login.GetAPIHost()
	request.URL.Path = strings.ReplaceAll(request.URL.Path, ":userID", string(login.UserID))
	request.URL.Path = strings.ReplaceAll(request.URL.Path, ":organizationID", string(login.OrgID))

	return c.request(request, target)
}

func (c *client) requestWithAgentAuthentication(request *http.Request, target any) error {
	login, err := c.getAgentLogin(request.Context())
	if err != nil {
		return err
	}

	request.URL.Scheme = "https"
	request.URL.Host = login.GetAPIHost()
	request.URL.Path = strings.ReplaceAll(request.URL.Path, ":userID", string(login.UserID))
	request.URL.Path = strings.ReplaceAll(request.URL.Path, ":organizationID", string(login.OrgID))

	return c.request(request, target)
}

func (c *client) getLogin(
	ctx context.Context,
	apiContextPath string,
	loginMutex *sync.Mutex,
	loginResponse *loginResponse,
) (*loginResponse, error) {
	{ // check for existing login
		if loginResponse != nil {
			return loginResponse, nil
		}

		loginMutex.Lock()
		defer loginMutex.Unlock()

		if loginResponse != nil {
			return loginResponse, nil
		}
	}

	login, err := c.endpointLogin(ctx, apiContextPath)
	if err != nil {
		return nil, err
	}

	loginResponse = &login

	if err := c.endpointGetSessionMetadata(ctx, apiContextPath); err != nil {
		return nil, err
	}

	loginState, err := c.endpointGetLoginState(ctx, apiContextPath)
	if err != nil {
		return nil, err
	}

	if loginState == "SELECT_STATION" {
		if err := c.endpointStartSession(ctx, apiContextPath); err != nil {
			return nil, err
		}
	}

	return c.supervisorLogin, nil
}

func (c *client) getSupervisorLogin(ctx context.Context) (*loginResponse, error) {
	return c.getLogin(ctx, "supsvcs/rs/svc", c.supervisorLoginMutex, c.supervisorLogin)
}

func (c *client) getAgentLogin(ctx context.Context) (*loginResponse, error) {
	return c.getLogin(ctx, "appsvcs/rs/svc", c.agentLoginMutex, c.agentLogin)
}

func (c *client) endpointLogin(ctx context.Context, apiContextPath string) (loginResponse, error) {
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
		fmt.Sprintf("https://app.five9.com/%s/auth/login", apiContextPath),
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

func (c *client) endpointGetSessionMetadata(ctx context.Context, apiContextPath string) error {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("/%s/auth/metadata", apiContextPath),
		http.NoBody,
	)
	if err != nil {
		return err
	}

	return c.requestWithSupervisorAuthentication(request, nil)
}

func (c *client) endpointGetLoginState(ctx context.Context, apiContextPath string) (userLoginState, error) {
	switch apiContextPath {
	case agentAPIContextPath:
		apiContextPath = fmt.Sprintf("%s/agents", agentAPIContextPath)
	case supervisorAPIContextPath:
		apiContextPath = fmt.Sprintf("%s/supervisors", supervisorAPIContextPath)
	default:
		return "", errors.New("unknown API Context Path")
	}

	var target userLoginState

	tries := 0
	for tries < 3 {
		tries++

		request, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			fmt.Sprintf("/%s/:userID/login_state", apiContextPath),
			http.NoBody,
		)
		if err != nil {
			return "", err
		}

		if err := c.requestWithSupervisorAuthentication(request, &target); err != nil {
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

func (c *client) endpointStartSession(ctx context.Context, apiContextPath string) error {
	switch apiContextPath {
	case agentAPIContextPath:
		apiContextPath = fmt.Sprintf("%s/agents", agentAPIContextPath)
	case supervisorAPIContextPath:
		apiContextPath = fmt.Sprintf("%s/supervisors", supervisorAPIContextPath)
	default:
		return errors.New("unknown API Context Path")
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		fmt.Sprintf("/%s/:userID/session_start?force=true", apiContextPath),
		structToReaderCloser(StationInfo{
			StationID:   "",
			StationType: "EMPTY",
		}),
	)
	if err != nil {
		return err
	}

	if err := c.requestWithSupervisorAuthentication(request, nil); err != nil {
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
