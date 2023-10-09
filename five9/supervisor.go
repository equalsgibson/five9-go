package five9

import (
	"context"
	"net/http"
)

type SupervisorService struct {
	client           *client
	websocketHandler websocketHandler
	cache            *supervisorWebsocketCache
	websocketReady   chan bool
}

func (s *SupervisorService) getAllDomainUsers(ctx context.Context) ([]AgentInfo, error) {
	var target []AgentInfo

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"/supsvcs/rs/svc/orgs/:organizationID/users",
		http.NoBody,
	)
	if err != nil {
		return nil, err
	}

	if err := s.client.requestWithAuthentication(request, &target); err != nil {
		return nil, err
	}

	return target, nil
}

func (s *SupervisorService) requestWebSocketFullStatistics(ctx context.Context) error {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		"/supsvcs/rs/svc/supervisors/:userID/request_full_statistics",
		http.NoBody,
	)
	if err != nil {
		return err
	}

	if err := s.client.requestWithAuthentication(request, nil); err != nil {
		return err
	}

	return nil
}
