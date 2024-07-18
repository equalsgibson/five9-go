package five9

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/equalsgibson/five9-go/five9/five9types"
)

type client struct {
	httpClient           *http.Client
	credentials          five9types.PasswordCredentials
	requestPreProcessors []func(r *http.Request) error
}

const (
	supervisorAPIContextPath = "supsvcs/rs/svc"
	agentAPIContextPath      = "appsvcs/rs/svc"
	statisticsAPIContextPath = "strsvcs/rs/svc"
	agentAPIPath             = "agents"
	supervisorAPIPath        = "supervisors"
	// statisticsPath = ""
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

func structToReaderCloser(v any) io.Reader {
	vBytes, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	return bytes.NewReader(vBytes)
}
