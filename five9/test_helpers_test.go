package five9_test

import (
	"io"
	"os"
	"testing"
)

func createIoReadCloserFromFile(t *testing.T, filePath string) io.ReadCloser {
	t.Helper()

	file, err := os.Open(filePath)
	if err != nil {
		t.Fatalf("Response Body File Not Found: %s", filePath)
	}

	return io.NopCloser(file)
}

// func createByteSliceFromFile(t *testing.T, filePath string) []byte {
// 	t.Helper()

// 	fileBytes, err := os.ReadFile(filePath)
// 	if err != nil {
// 		t.Fatalf("File Not Found: %s", filePath)
// 	}

// 	return fileBytes
// }

// The below requests run in order when first starting the websocket service.
// func generateWSLoginRequestFuncs(t *testing.T) []func(r *http.Request) (*http.Response, error) {
// 	t.Helper()

// 	return []func(r *http.Request) (*http.Response, error){
// 		func(r *http.Request) (*http.Response, error) { // https://app.five9.com/supsvcs/rs/svc/auth/login
// 			return &http.Response{
// 				Body:       createIoReadCloserFromFile(t, "test/supervisorLogin_200.json"),
// 				StatusCode: http.StatusOK,
// 			}, nil
// 		},
// 		func(r *http.Request) (*http.Response, error) { // supsvcs/rs/svc/auth/metadata
// 			return &http.Response{
// 				Body:       createIoReadCloserFromFile(t, "test/auth_metadata_200.json"),
// 				StatusCode: http.StatusOK,
// 			}, nil
// 		},
// 		func(r *http.Request) (*http.Response, error) { // supsvcs/rs/svc/supervisors/:userID/login_state
// 			return &http.Response{
// 				Body:       createIoReadCloserFromFile(t, "test/loginState_selectStation_200.json"),
// 				StatusCode: http.StatusOK,
// 			}, nil
// 		},
// 		func(r *http.Request) (*http.Response, error) { // supsvcs/rs/svc/supervisors/:userID/session_start?force=true
// 			return &http.Response{
// 				Body:       http.NoBody,
// 				StatusCode: http.StatusNoContent,
// 			}, nil
// 		},
// 		func(r *http.Request) (*http.Response, error) { // supsvcs/rs/svc/orgs/:organizationID/users
// 			return &http.Response{
// 				Body:       createIoReadCloserFromFile(t, "test/supervisor_getAllUsers_200.json"),
// 				StatusCode: http.StatusOK,
// 			}, nil
// 		},
// 		func(r *http.Request) (*http.Response, error) { // request_full_statistics
// 			return &http.Response{
// 				Body:       http.NoBody,
// 				StatusCode: http.StatusNoContent,
// 			}, nil
// 		},
// 	}
// }
