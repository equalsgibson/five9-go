package five9

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/equalsgibson/five9-go/five9/five9types"
)

type authenticationState struct {
	client         *client
	loginResponse  *five9types.LoginResponse
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

	var latestAttemptErr error
	tries := 0
	for tries < 3 {
		tries++
		latestAttemptErr = a.client.request(request, target)
		if latestAttemptErr != nil {
			if five9Error, ok := latestAttemptErr.(*Error); ok {
				if five9Error.StatusCode == http.StatusUnauthorized {
					// The login is not registered by other endpoints for a short time.
					// I think this has to do with Five9 propagating the session across their data centers.
					// We login using the app.five9.com domain but then make subsequent calls to the data center specific domain
					time.Sleep(time.Second * 2)

					continue
				}

				// Five9 reply with Status 435 if a service has been migrated. This is not an official status code, so check directly.
				if five9Error.StatusCode == int(435) {
					// Clear out the login state
					a.loginMutex.Lock()
					defer a.loginMutex.Unlock()

					a.loginResponse = nil

					return latestAttemptErr
				}
			}
		}

		return nil
	}

	return latestAttemptErr
}
func (a *authenticationState) requestDownloadWithAuthentication(request *http.Request) (*http.Response, error) {
	login, err := a.getLogin(request.Context())
	if err != nil {
		return nil, err
	}

	request.URL.Scheme = "https"
	request.URL.Host = login.GetAPIHost()
	request.URL.Path = strings.ReplaceAll(request.URL.Path, ":userID", string(login.UserID))
	request.URL.Path = strings.ReplaceAll(request.URL.Path, ":organizationID", string(login.OrgID))

	var latestAttemptErr error
	tries := 0
	for tries < 3 {
		tries++
		response, latestAttemptErr := a.client.requestDownload(request)
		if latestAttemptErr != nil {
			if five9Error, ok := latestAttemptErr.(*Error); ok {
				if five9Error.StatusCode == http.StatusUnauthorized {
					time.Sleep(time.Second * 2)
					continue
				}

				if five9Error.StatusCode == int(435) {
					a.loginMutex.Lock()
					defer a.loginMutex.Unlock()

					a.loginResponse = nil

					return nil, latestAttemptErr
				}
			}
		} else {
			return response, nil
		}
	}

	return nil, latestAttemptErr
}

func (a *authenticationState) getLogin(
	ctx context.Context,
) (*five9types.LoginResponse, error) {
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

	switch loginState {
	case five9types.UserLoginStateSelectStation: // Standard response after logging in.
		if err := a.endpointStartSession(ctx); err != nil {
			return nil, err
		}

		// Check the login state after starting the session
		newLoginState, err := a.endpointGetLoginState(ctx)
		if err != nil {
			return nil, err
		}

		if newLoginState == five9types.UserLoginStateAcceptNotice {
			if err := a.handleMaintenanceNotices(ctx); err != nil {
				return nil, err
			}
		}

	case five9types.UserLoginStateAcceptNotice: // Can occur if Five9 have issued a maintenance notice
		if err := a.handleMaintenanceNotices(ctx); err != nil {
			return nil, err
		}

	case five9types.UserLoginStateRelogin: // Can occur if the service has been migrated
		if err := a.endpointRestartSession(ctx); err != nil {
			return nil, err
		}

		// Check the login state after restarting the session
		newLoginState, err := a.endpointGetLoginState(ctx)
		if err != nil {
			return nil, err
		}

		if newLoginState == five9types.UserLoginStateAcceptNotice {
			if err := a.handleMaintenanceNotices(ctx); err != nil {
				return nil, err
			}
		}
	}

	return a.loginResponse, nil
}

func (a *authenticationState) endpointLogin(ctx context.Context) (five9types.LoginResponse, error) {
	payload := five9types.LoginPayload{
		PasswordCredentials: a.client.credentials,
		AppKey:              "web-ui",
		Policy:              five9types.PolicyForceIn,
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("https://app.five9.com/%s/auth/login", a.apiContextPath),
		structToReaderCloser(payload),
	)
	if err != nil {
		return five9types.LoginResponse{}, err
	}

	target := five9types.LoginResponse{}

	if err := a.client.request(request, &target); err != nil {
		return five9types.LoginResponse{}, err
	}

	return target, nil
}

func (a *authenticationState) endpointGetLoginState(ctx context.Context) (five9types.UserLoginState, error) {
	path := ""

	switch a.apiContextPath {
	case supervisorAPIContextPath:
		path = supervisorAPIPath
	case agentAPIContextPath:
		path = agentAPIPath
	case statisticsAPIContextPath:
		path = supervisorAPIPath
	}

	if path == "" {
		return "", errors.New("could not set API context path, library error")
	}

	var target five9types.UserLoginState

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
		return "", err
	}

	return target, nil
}

func (a *authenticationState) endpointStartSession(ctx context.Context) error {
	path := agentAPIPath
	if a.apiContextPath == supervisorAPIContextPath {
		path = supervisorAPIPath
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		fmt.Sprintf(
			"/%s/%s/:userID/session_start?force=true",
			a.apiContextPath,
			path,
		),
		structToReaderCloser(five9types.StationInfo{
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

func (a *authenticationState) endpointRestartSession(ctx context.Context) error {
	path := agentAPIPath
	if a.apiContextPath == supervisorAPIContextPath {
		path = supervisorAPIPath
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		fmt.Sprintf(
			"/%s/%s/:userID/session_restart",
			a.apiContextPath,
			path,
		),
		http.NoBody,
	)
	if err != nil {
		return err
	}

	if err := a.requestWithAuthentication(request, nil); err != nil {
		return err
	}

	return nil
}

func (a *authenticationState) handleMaintenanceNotices(ctx context.Context) error {
	notices, err := a.endpointGetMaintenanceNotices(ctx)
	if err != nil {
		return err
	}

	for _, notice := range notices {
		if err := a.endpointAcceptMaintenanceNotice(ctx, notice.ID); err != nil {
			return err
		}
	}

	loginState, err := a.endpointGetLoginState(ctx)
	if err != nil {
		return err
	}

	if loginState != five9types.UserLoginStateWorking {
		if err := a.endpointStartSession(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *authenticationState) endpointGetMaintenanceNotices(ctx context.Context) ([]five9types.MaintenanceNoticeInfo, error) {
	path := agentAPIPath
	if a.apiContextPath == supervisorAPIContextPath {
		path = supervisorAPIPath
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf(
			"/%s/%s/:userID/maintenance_notices",
			a.apiContextPath,
			path,
		),
		http.NoBody,
	)
	if err != nil {
		return nil, err
	}

	target := []five9types.MaintenanceNoticeInfo{}

	if err := a.requestWithAuthentication(request, &target); err != nil {
		return nil, err
	}

	return target, nil
}

func (a *authenticationState) endpointAcceptMaintenanceNotice(
	ctx context.Context,
	maintenanceNoticeID five9types.MaintenanceNoticeID,
) error {
	path := agentAPIPath
	if a.apiContextPath == supervisorAPIContextPath {
		path = supervisorAPIPath
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		fmt.Sprintf(
			"/%s/%s/:userID/maintenance_notices/%s/accept",
			a.apiContextPath,
			path,
			maintenanceNoticeID,
		),
		http.NoBody,
	)
	if err != nil {
		return err
	}

	target := five9types.MaintenanceNoticeInfo{}

	return a.requestWithAuthentication(request, &target)
}
