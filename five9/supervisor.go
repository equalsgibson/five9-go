package five9

import (
	"context"
	"net/http"
)

type SupervisorService struct {
	client              *client
	websocketHandler    websocketHandler
	webSocketCache      *supervisorWebsocketCache
	domainMetadataCache *domainMetadata
	websocketReady      chan bool
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

func (s *SupervisorService) getAllReasonCodes(ctx context.Context) ([]ReasonCodeInfo, error) {
	logoutCodes, err := s.getAllLogoutReasonCodes(ctx)
	if err != nil {
		return nil, err
	}

	notReadyCodes, err := s.getAllNotReadyReasonCodes(ctx)
	if err != nil {
		return nil, err
	}

	return append(logoutCodes, notReadyCodes...), nil
}

func (s *SupervisorService) getAllLogoutReasonCodes(ctx context.Context) ([]ReasonCodeInfo, error) {
	var target []ReasonCodeInfo

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"/appsvcs/rs/svc/orgs/:organizationID/logout_reason_codes",
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

func (s *SupervisorService) getAllNotReadyReasonCodes(ctx context.Context) ([]ReasonCodeInfo, error) {
	var target []ReasonCodeInfo

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"/appsvcs/rs/svc/orgs/:organizationID/not_ready_reason_codes",
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

func (s *SupervisorService) DomainReasonCodes(ctx context.Context) map[ReasonCodeID]ReasonCodeInfo {
	return s.domainMetadataCache.reasonCodes
}
