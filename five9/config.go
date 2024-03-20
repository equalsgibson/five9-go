package five9

import (
	"net/http"

	"github.com/equalsgibson/five9-go/five9/five9types"
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

// func SetMaxCacheAge(maxAge *time.Duration) ConfigFunc {
// 	return func(s *Service) {
// 		s.agentService.authState.client.httpClient.Transport = roundTripper
// 	}
// }

func SetAgentStateCache(
	cache Cache[five9types.UserID, five9types.AgentState],
) ConfigFunc {
	return func(s *Service) {
		s.supervisorService.webSocketCache.agentState = cache
	}
}

func AcdStateCache(
	cache Cache[five9types.QueueID, five9types.ACDState],
) ConfigFunc {
	return func(s *Service) {
		s.supervisorService.webSocketCache.acdState = cache
	}
}
