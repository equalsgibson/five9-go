package five9

import (
	"context"
	"net/http"
	"time"
)

type AgentService struct {
	authState           *authenticationState
	websocketHandler    websocketHandler
	webSocketCache      *agentWebsocketCache
	domainMetadataCache *domainMetadata
	websocketReady      chan bool
}

type agentWebsocketCache struct {
	lastPong *time.Time
}

func (s *AgentService) GetAllReasonCodes(ctx context.Context) ([]ReasonCodeInfo, error) {
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

func (s *AgentService) getAllLogoutReasonCodes(ctx context.Context) ([]ReasonCodeInfo, error) {
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

	if err := s.authState.requestWithAuthentication(request, &target); err != nil {
		return nil, err
	}

	return target, nil
}

func (s *AgentService) getAllNotReadyReasonCodes(ctx context.Context) ([]ReasonCodeInfo, error) {
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

	if err := s.authState.requestWithAuthentication(request, &target); err != nil {
		return nil, err
	}

	return target, nil
}
