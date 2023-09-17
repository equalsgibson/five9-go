package five9

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

type client struct {
	httpClient          *http.Client
	authenticationMutex *sync.Mutex
	credentials         PasswordCredentials
	restAPIURL          string
	context             string
	token               *AuthenticationTokenID
	farmID              *FarmID
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

func (c *client) json(ctx context.Context, method string, path string, body io.Reader, target any) error {
	request, err := http.NewRequestWithContext(ctx, method, path, body)
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

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
	// url, err := url.Parse("https://app.five9.com/")
	// if err != nil {
	// 	return err
	// }

	// cookies := c.httpClient.Jar.Cookies(url)

	if c.token != nil {
		return nil
	}

	c.authenticationMutex.Lock()
	defer c.authenticationMutex.Unlock()

	if c.token != nil {
		return nil
	}

	target := LoginSupervisorResponse{}
	payload, err := json.Marshal(LoginPayload{
		PasswordCredentials: c.credentials,
		AppKey:              "web-ui",
		Policy:              PolicyAttachExisting,
	})
	if err != nil {
		panic(err)
	}

	if err := c.json(
		ctx,
		http.MethodPost,
		fmt.Sprintf("https://app.five9.com/%s/auth/login", c.context),
		bytes.NewReader(payload),
		&target,
	); err != nil {
		return err
	}

	// u, err := url.Parse("https://app.five9.com/")
	// if err != nil {
	// 	return err
	// }

	c.farmID = &target.Context.FarmID
	c.token = &target.TokenID

	target2 := SupervisorMetadataResponse{}

	if err := c.json(
		ctx,
		http.MethodGet,
		fmt.Sprintf("https://%s/%s/auth/metadata", target.Metadata.DataCenters[0].APIURLs[0].Host, c.context),
		nil,
		&target2,
	); err != nil {
		return err
	}

	return nil
}
