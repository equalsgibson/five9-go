package five9

import "net/http"

type PasswordCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ConfigFunc func(*Service)

func AddRequestPreprocessor(things ...func(*http.Request) error) ConfigFunc {
	return func(s *Service) {
		s.client.requestPreProcessors = append(
			s.client.requestPreProcessors,
			things...,
		)
	}
}

func setWebsocketHandler(w websocketHandler) ConfigFunc {
	return func(s *Service) {
		s.supervisorService.websocketHandler = w
	}
}

func setRoundTripper(roundTripper http.RoundTripper) ConfigFunc {
	return func(s *Service) {
		s.client.httpClient.Transport = roundTripper
	}
}
