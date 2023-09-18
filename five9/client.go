package five9

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
)

const Five9APILoginURL string = "app.five9.com"

type client struct {
	httpClient          *http.Client
	authenticationMutex *sync.Mutex
	credentials         PasswordCredentials
	restAPIURL          string
	context             string
	authenticated       bool
}

func (c *client) do(r *http.Request) (*http.Response, error) {
	r.URL.Scheme = "https"
	r.Header.Set("Accept", "application/json")

	response, err := c.httpClient.Do(r)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *client) json(ctx context.Context, method string, path string, body io.Reader, target any, requestPreProcessors ...RequestPreProcessorFunc) error {
	request, err := http.NewRequestWithContext(ctx, method, path, body)
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	for _, preprocessor := range requestPreProcessors {
		preprocessor(request)
	}

	response, err := c.do(request)
	if err != nil {
		return err
	}

	if target != nil {
		defer response.Body.Close()

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

func (c *client) authenticate(ctx context.Context) error {
	if c.authenticated {
		return nil
	}

	c.authenticationMutex.Lock()
	defer c.authenticationMutex.Unlock()

	if c.authenticated {
		return nil
	}

	// Reset the REST API URL to empty
	c.restAPIURL = ""

	five9LoginURL := url.URL{
		Scheme: "https",
		Host:   Five9APILoginURL,
		Path:   fmt.Sprintf("%s/auth/login", c.context),
	}

	payload, err := json.Marshal(
		LoginPayload{
			PasswordCredentials: c.credentials,
			AppKey:              "web-ui",
			Policy:              PolicyAttachExisting,
		},
	)

	if err != nil {
		panic(err)
	}

	target := LoginSupervisorResponse{}
	if err := c.json(
		ctx,
		http.MethodPost,
		five9LoginURL.String(),
		bytes.NewReader(payload),
		&target,
	); err != nil {
		return err
	}

	// Loop over the Datacenters to find the active datacenter, then set the REST API URL.
	for _, datacenter := range target.Metadata.DataCenters {
		if !datacenter.Active {
			continue
		}

		if len(datacenter.APIURLs) > 0 {
			c.restAPIURL = datacenter.APIURLs[0].Host
		}
		break
	}

	if c.restAPIURL == "" {
		return errors.New("unable to get REST API URL from login response")
	}

	five9BaseURL := url.URL{
		Scheme: "https",
		Host:   Five9APILoginURL,
	}

	// Confirm the cookies are set correctly
	cookies := c.httpClient.Jar.Cookies(&five9BaseURL)
	farmSet := false
	authSet := false

	for _, cookie := range cookies {
		if farmSet && authSet {
			break
		}
		if cookie.Name == "farmId" {
			farmSet = true
			continue
		}
		if cookie.Name == "Authorization" {
			authSet = true
			continue
		}
	}

	if !farmSet || !authSet {
		return errors.New("unable to get authorization and farm cookies")
	}

	return nil
}

type RequestPreProcessor interface {
	ProcessRequest(*http.Request) error
}

type RequestPreProcessorFunc func(*http.Request) error

func (p RequestPreProcessorFunc) ProcessRequest(r *http.Request) error {
	return p(r)
}
