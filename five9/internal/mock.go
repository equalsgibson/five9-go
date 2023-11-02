package internal

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/equalsgibson/five9-go/five9/five9types"
)

type mockFive9Server struct {
	Routes map[string]http.Handler
	Users  map[five9types.UserID]*five9types.UserLoginState
}

func (m *mockFive9Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	matchedRoute, ok := m.Routes[r.URL.Path]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	matchedRoute.ServeHTTP(w, r)
}

func NewMockFive9Server(t *testing.T, configFuncs ...ConfigFunc) *httptest.Server {
	defaultRoutes := map[string]http.Handler{}

	mockFive9Server := mockFive9Server{
		Routes: defaultRoutes,
	}

	server := httptest.NewTLSServer(&mockFive9Server)

	defaultRoutes["/supsvcs/rs/svc/auth/login"] = mockFive9Server.handleLogin(strings.TrimPrefix(server.URL, "https://"))
	defaultRoutes["/supsvcs/rs/svc/auth/metadata"] = mockFive9Server.handleMetadata(strings.TrimPrefix(server.URL, "https://"))
	defaultRoutes["/supsvcs/rs/svc/supervisors/"] = mockFive9Server.handleGetLoginState()

	for _, configFunc := range configFuncs {
		configFunc(&mockFive9Server)
	}

	return server
}

func (m *mockFive9Server) handleLogin(url string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urlInfo := strings.Split(url, ":")

		response := &five9types.LoginResponse{}
		file := createByteSliceFromFile("mockResponses/supervisorLogin_200.json")

		if err := json.Unmarshal(file, response); err != nil {
			log.Printf("ERROR! -> %s", err.Error())
			return
		}

		response.Metadata.DataCenters = []five9types.DataCenter{
			{
				Name:  "Test",
				UI:    nil,
				Login: nil,
				API: []five9types.Server{
					{
						Host:     urlInfo[0],
						Port:     urlInfo[1],
						RouteKey: "FAKEKEY456",
						Version:  "13.0.183",
					},
				},
				Active: true,
			},
		}

		b, _ := json.Marshal(response)

		w.WriteHeader(http.StatusOK)
		w.Write(b)
	})
}

func (m *mockFive9Server) handleMetadata(url string) http.Handler {
	return m.handleLogin(url)
}

func (m *mockFive9Server) handleGetLoginState() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("\"WORKING\""))
	})
}

type ConfigFunc func(*mockFive9Server)

func SetRoutes(routes map[string]http.Handler) ConfigFunc {
	return func(s *mockFive9Server) {
		for route, handler := range routes {
			s.Routes[route] = handler
		}
	}
}

func createByteSliceFromFile(filePath string) []byte {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil
	}

	return fileBytes
}
