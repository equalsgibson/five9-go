package five9

import "net/http"

type ConfigFunc func(*Service)

func AddRequestPreprocessor(things ...func(*http.Request) error) ConfigFunc {
	return func(s *Service) {
		s.agentService.authState.client.requestPreProcessors = append(
			s.agentService.authState.client.requestPreProcessors,
			things...,
		)
	}
}

func SetWebsocketHandler(w webSocketHandler) ConfigFunc {
	return func(s *Service) {
		s.supervisorService.webSocketHandler = w
	}
}

func SetRoundTripper(roundTripper http.RoundTripper) ConfigFunc {
	return func(s *Service) {
		s.agentService.authState.client.httpClient.Transport = roundTripper
	}
}

func SetFive9ServerLoginURL(url string) ConfigFunc {
	return func(s *Service) {
		s.agentService.authState.loginURL = url
		s.supervisorService.authState.loginURL = url
	}
}
