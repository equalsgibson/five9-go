package five9_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/equalsgibson/five9-go/five9"
	"github.com/equalsgibson/five9-go/five9/five9types"
)

func Test_Authentication_Success(t *testing.T) {
	ctx := context.Background()
	madeAllExpectedAPICalls := false

	mockRoundTripper := MockRoundTripper{
		Func: []func(r *http.Request) (*http.Response, error){
			func(r *http.Request) (*http.Response, error) { // https://app.five9.com/supsvcs/rs/svc/auth/login
				return &http.Response{
					Body:       createIoReadCloserFromFile(t, "test/supervisorLogin_200.json"),
					StatusCode: http.StatusOK,
				}, nil
			},
			func(r *http.Request) (*http.Response, error) { // supsvcs/rs/svc/auth/metadata
				return &http.Response{
					Body:       createIoReadCloserFromFile(t, "test/auth_metadata_200.json"),
					StatusCode: http.StatusOK,
				}, nil
			},
			func(r *http.Request) (*http.Response, error) { // supsvcs/rs/svc/supervisors/:userID/login_state
				return &http.Response{
					Body:       createIoReadCloserFromFile(t, "test/loginState_selectStation_200.json"),
					StatusCode: http.StatusOK,
				}, nil
			},
			func(r *http.Request) (*http.Response, error) { // supsvcs/rs/svc/supervisors/:userID/session_start?force=true
				return &http.Response{
					Body:       http.NoBody,
					StatusCode: http.StatusNoContent,
				}, nil
			},
			func(r *http.Request) (*http.Response, error) { // supsvcs/rs/svc/orgs/:organizationID/users
				madeAllExpectedAPICalls = true

				return &http.Response{
					Body:       createIoReadCloserFromFile(t, "test/supervisor_getAllUsers_200.json"),
					StatusCode: http.StatusOK,
				}, nil
			},
		},
	}

	s := five9.NewService(
		five9types.PasswordCredentials{},
		five9.SetRoundTripper(&mockRoundTripper),
		five9.AddRequestPreprocessor(func(r *http.Request) error {
			t.Logf("API Call Made: [%s] %s\n", r.Method, r.URL.String())

			return nil
		}),
	)

	// This could be any request that requires authentication
	_, err := s.Supervisor().GetAllDomainUsers(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if !madeAllExpectedAPICalls {
		t.Fatalf("did not make all expected API calls - %d api requests remaining in queue", len(mockRoundTripper.Func))
	}
}

func Test_Authentication_Reuse_LoginState_Success(t *testing.T) {
	ctx := context.Background()
	madeAllExpectedAPICalls := false

	mockRoundTripper := MockRoundTripper{
		Func: []func(r *http.Request) (*http.Response, error){
			func(r *http.Request) (*http.Response, error) { // https://app.five9.com/supsvcs/rs/svc/auth/login
				return &http.Response{
					Body:       createIoReadCloserFromFile(t, "test/supervisorLogin_200.json"),
					StatusCode: http.StatusOK,
				}, nil
			},
			func(r *http.Request) (*http.Response, error) { // supsvcs/rs/svc/auth/metadata
				return &http.Response{
					Body:       createIoReadCloserFromFile(t, "test/auth_metadata_200.json"),
					StatusCode: http.StatusOK,
				}, nil
			},
			func(r *http.Request) (*http.Response, error) { // supsvcs/rs/svc/supervisors/:userID/login_state
				return &http.Response{
					Body:       createIoReadCloserFromFile(t, "test/loginState_selectStation_200.json"),
					StatusCode: http.StatusOK,
				}, nil
			},
			func(r *http.Request) (*http.Response, error) { // supsvcs/rs/svc/supervisors/:userID/session_start?force=true
				return &http.Response{
					Body:       http.NoBody,
					StatusCode: http.StatusNoContent,
				}, nil
			},
			func(r *http.Request) (*http.Response, error) { // supsvcs/rs/svc/orgs/:organizationID/users
				return &http.Response{
					Body:       createIoReadCloserFromFile(t, "test/supervisor_getAllUsers_200.json"),
					StatusCode: http.StatusOK,
				}, nil
			},
			func(r *http.Request) (*http.Response, error) { // supsvcs/rs/svc/orgs/:organizationID/users
				madeAllExpectedAPICalls = true

				return &http.Response{
					Body:       createIoReadCloserFromFile(t, "test/supervisor_getAllUsers_200.json"),
					StatusCode: http.StatusOK,
				}, nil
			},
		},
	}

	s := five9.NewService(
		five9types.PasswordCredentials{},
		five9.SetRoundTripper(&mockRoundTripper),
		five9.AddRequestPreprocessor(func(r *http.Request) error {
			t.Logf("API Call Made: [%s] %s\n", r.Method, r.URL.String())

			return nil
		}),
	)

	// This could be any request that requires authentication
	_, err := s.Supervisor().GetAllDomainUsers(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Call the request again to ensure we are reusing login state after it is set
	_, err = s.Supervisor().GetAllDomainUsers(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if !madeAllExpectedAPICalls {
		t.Fatalf("did not make all expected API calls - %d api requests remaining in queue", len(mockRoundTripper.Func))
	}
}

func Test_Authentication_AcceptNotices(t *testing.T) {
	ctx := context.Background()
	madeAllExpectedAPICalls := false

	mockRoundTripper := MockRoundTripper{
		Func: []func(r *http.Request) (*http.Response, error){
			func(r *http.Request) (*http.Response, error) { // https://app.five9.com/supsvcs/rs/svc/auth/login
				return &http.Response{
					Body:       createIoReadCloserFromFile(t, "test/supervisorLogin_200.json"),
					StatusCode: http.StatusOK,
				}, nil
			},
			func(r *http.Request) (*http.Response, error) { // supsvcs/rs/svc/auth/metadata
				return &http.Response{
					Body:       createIoReadCloserFromFile(t, "test/auth_metadata_200.json"),
					StatusCode: http.StatusOK,
				}, nil
			},
			func(r *http.Request) (*http.Response, error) { // supsvcs/rs/svc/supervisors/:userID/login_state
				return &http.Response{
					Body:       createIoReadCloserFromFile(t, "test/loginState_acceptNotice_200.json"),
					StatusCode: http.StatusOK,
				}, nil
			},
			func(r *http.Request) (*http.Response, error) { // supsvcs/rs/svc/supervisors/:userID/maintenance_notices
				return &http.Response{
					Body:       createIoReadCloserFromFile(t, "test/getMaintenanceNotices_Notices_200.json"),
					StatusCode: http.StatusOK,
				}, nil
			},
			func(r *http.Request) (*http.Response, error) { // supsvcs/rs/svc/supervisors/:userID/maintenance_notices/%s/accept
				return &http.Response{
					Body:       createIoReadCloserFromFile(t, "test/acceptMaintenanceNotice_8213_200.json"),
					StatusCode: http.StatusOK,
				}, nil
			},
			func(r *http.Request) (*http.Response, error) { // supsvcs/rs/svc/supervisors/:userID/maintenance_notices/%s/accept
				return &http.Response{
					Body:       createIoReadCloserFromFile(t, "test/acceptMaintenanceNotice_8214_200.json"),
					StatusCode: http.StatusOK,
				}, nil
			},
			func(r *http.Request) (*http.Response, error) { // supsvcs/rs/svc/supervisors/:userID/session_start?force=true
				return &http.Response{
					Body:       http.NoBody,
					StatusCode: http.StatusNoContent,
				}, nil
			},
			func(r *http.Request) (*http.Response, error) { // supsvcs/rs/svc/orgs/:organizationID/users
				madeAllExpectedAPICalls = true

				return &http.Response{
					Body:       createIoReadCloserFromFile(t, "test/supervisor_getAllUsers_200.json"),
					StatusCode: http.StatusOK,
				}, nil
			},
		},
	}

	s := five9.NewService(
		five9types.PasswordCredentials{},
		five9.SetRoundTripper(&mockRoundTripper),
		five9.AddRequestPreprocessor(func(r *http.Request) error {
			t.Logf("API Call Made: [%s] %s\n", r.Method, r.URL.String())

			return nil
		}),
	)
	// This could be any request that requires authentication
	_, err := s.Supervisor().GetAllDomainUsers(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if !madeAllExpectedAPICalls {
		t.Fatalf("did not make all expected API calls - %d api requests remaining in queue", len(mockRoundTripper.Func))
	}
}

func TestThis(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	}))

	testService := five9.NewService(five9types.PasswordCredentials{}, five9.SetFive9ServerLoginURL(server.URL))

}
