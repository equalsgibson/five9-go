package five9

import (
	"crypto/tls"
	"net/http"
)

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

func SetTLSInsecureSkipVerify() ConfigFunc {
	return func(s *Service) {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
		}
		transport := http.Transport{
			TLSClientConfig: tlsConfig,
		}

		s.agentService.authState.client.httpClient.Transport = &transport
	}
}

func SetFive9ServerLoginURL(url string) ConfigFunc {
	return func(s *Service) {
		s.agentService.authState.loginURL = url
		s.supervisorService.authState.loginURL = url
	}
}
